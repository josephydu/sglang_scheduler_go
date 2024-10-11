package controller

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"

	"github.com/go-resty/resty/v2"
	"github.com/pebbe/zmq4"

	"sglang_scheduler_go/models"
	"sglang_scheduler_go/server_args"
)

type LoadBalanceMethod int

const (
	RoundRobin LoadBalanceMethod = iota
	ShortestQueue
	ResourcesAware
	PowerOf2Choice
	PreRadix
	MultiTurn
	Bucket
)

func LoadBalanceMethodFromString(method string) (LoadBalanceMethod, error) {
	switch method {
	case "round_robin":
		return RoundRobin, nil
	case "shortest_queue":
		return ShortestQueue, nil
	case "resources_aware":
		return ResourcesAware, nil
	case "power_of_2_choice":
		return PowerOf2Choice, nil
	case "pre_radix":
		return PreRadix, nil
	case "multi_turn":
		return MultiTurn, nil
	case "bucket":
		return Bucket, nil
	default:
		return -1, fmt.Errorf("%s is not a valid load balance method", method)
	}
}

type Controller struct {
	serverArgs            *server_args.ServerArgs
	loadBalanceMethod     LoadBalanceMethod
	NodeList              []models.NodeInfo
	recvControllerProcess []*zmq4.Socket
	controllerInfoDict    map[string][3]int
	roundRobinCounter     int
	Dispatching           func([]models.Request, string) <-chan []byte
	context               *zmq4.Context
	logStep               int
	mutex                 sync.Mutex
}

func NewController(serverArgs *server_args.ServerArgs) *Controller {
	loadBalanceMethod, err := LoadBalanceMethodFromString(serverArgs.LoadBalanceMethod)
	if err != nil {
		log.Fatalf("Invalid load balance method: %v", err)
	}

	context, err := zmq4.NewContext()
	if err != nil {
		log.Fatalf("Failed to create ZeroMQ context: %v", err)
	}

	c := &Controller{
		serverArgs:            serverArgs,
		loadBalanceMethod:     loadBalanceMethod,
		NodeList:              make([]models.NodeInfo, 0),
		recvControllerProcess: make([]*zmq4.Socket, 0),
		controllerInfoDict:    make(map[string][3]int),
		roundRobinCounter:     0,
		context:               context,
		logStep:               10,
	}

	switch loadBalanceMethod {
	case RoundRobin:
		c.Dispatching = c.roundRobinScheduler
	case PowerOf2Choice:
		c.Dispatching = c.powerOf2ChoiceScheduler
	default:
		log.Fatalf("Unsupported load balance method: %v", loadBalanceMethod)
	}
	return c
}

func (c *Controller) AddNewNode(nodeInfo models.NodeInfo) {
	log.Printf("%s:%d is registering on server...", nodeInfo.Ip, nodeInfo.Port)
	c.NodeList = append(c.NodeList, nodeInfo)
	if nodeInfo.ControllerInfoPort != 0 {
		recvControllerInfo, err := c.context.NewSocket(zmq4.PULL)
		if err != nil {
			log.Fatalf("Failed to create zeroMQ socket: %v", err)
		}
		err = recvControllerInfo.Connect(fmt.Sprintf("tcp://%s:%d", nodeInfo.Ip, nodeInfo.ControllerInfoPort))
		if err != nil {
			log.Fatalf("Failed to Connect zeroMQ: %v", err)
		}

		go c.recvControllerInfoLoop(recvControllerInfo)
		c.recvControllerProcess = append(c.recvControllerProcess, recvControllerInfo)
	}
}

func (c *Controller) recvControllerInfoLoop(recvControllerInfo *zmq4.Socket) {
	for {
		recvControllerInfoStr, err := recvControllerInfo.Recv(zmq4.DONTWAIT)
		if err != nil {
			if errors.Is(err, zmq4.Errno(zmq4.ETERM)) {
				return
			}
			continue
		}

		if recvControllerInfoStr != "" {
			var ip string
			var port, availableMemory, numRunning, numWaiting int
			_, err := fmt.Sscanf(recvControllerInfoStr, "%s,%d,%d,%d,%d", &ip, &port, &availableMemory, &numRunning, &numWaiting)
			if err != nil {
				continue
			}
			c.mutex.Lock()
			c.controllerInfoDict[fmt.Sprintf("%s:%d", ip, port)] = [3]int{availableMemory, numWaiting, numWaiting}
			c.mutex.Unlock()
		}
	}
}

func (c *Controller) roundRobinScheduler(inputRequests []models.Request, baseUrl string) <-chan []byte {
	out := make(chan []byte)
	go func() {
		defer close(out)
		if len(inputRequests) == 0 || len(c.NodeList) == 0 {
			return
		}

		client := resty.New()
		for _, req := range inputRequests {
			targetNode := c.NodeList[c.roundRobinCounter]
			c.roundRobinCounter = (c.roundRobinCounter + 1) % len(c.NodeList)

			payload := req.ToMap()
			url := fmt.Sprintf("http://%s:%d/%s", targetNode.Ip, targetNode.Port, baseUrl)
			resp, err := client.R().SetBody(payload).Post(url)

			if err != nil || resp.StatusCode() != 200 {
				log.Printf("Failed to retrieve data: %v", err)
				out <- []byte{}
				continue
			}
			out <- resp.Body()
		}
	}()
	return out
}

func (c *Controller) powerOf2ChoiceScheduler(inputRequests []models.Request, baseurl string) <-chan []byte {
	out := make(chan []byte)
	go func() {
		defer close(out)
		if len(inputRequests) == 0 || len(c.NodeList) == 0 || len(c.controllerInfoDict) == 0 {
			return
		}

		client := resty.New()
		for _, req := range inputRequests {
			var targetNode string
			if len(c.controllerInfoDict) == 1 {
				for k := range c.controllerInfoDict {
					targetNode = k
				}
			} else {
				keys := make([]string, 0, len(c.controllerInfoDict))

				for k := range c.controllerInfoDict {
					keys = append(keys, k)
				}

				k1, k2 := keys[rand.Intn(len(keys))], keys[rand.Intn(len(keys))]
				v1, v2 := c.controllerInfoDict[k1], c.controllerInfoDict[k2]

				if v1[2] != v2[2] {
					targetNode = k1
					if v2[2] < v2[2] {
						targetNode = k2
					}
				} else if v1[1] != v2[1] {
					targetNode = k1
					if v2[1] < v1[1] {
						targetNode = k2
					}
				} else if v1[0] != v2[0] {
					targetNode = k1
					if v2[0] > v1[0] {
						targetNode = k2
					}
				} else {
					targetNode = k1
				}
			}
			log.Printf("Chosen node in power of 2 choice is [%s]", targetNode)
			url := fmt.Sprintf("http://%s/%s", targetNode, baseurl)

			payload := req.ToMap()
			resp, err := client.R().SetBody(payload).Post(url)

			if err != nil || resp.StatusCode() != 200 {
				log.Printf("Failed to retrieve data: %v", err)
				out <- []byte{}
				continue
			}
			out <- resp.Body()
		}
	}()
	return out
}
