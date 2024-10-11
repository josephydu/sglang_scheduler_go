package main

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"sglang_scheduler_go/models"
)

func registerNodes(c *gin.Context) {
	var nodeInfo models.NodeInfo

	if err := c.ShouldBindJSON(&nodeInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if ctrl != nil {
		ctrl.AddNewNode(nodeInfo)
		c.JSON(http.StatusOK, gin.H{"message": "Register node to the controller successfully!"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Controller is not initialized."})
	}
}

func generate(c *gin.Context) {
	var req models.GenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if ctrl != nil {
		baseUrl := "generate"
		c.Stream(func(w io.Writer) bool {
			for data := range ctrl.Dispatching([]models.Request{req}, baseUrl) {
				if len(data) > 0 {
					w.Write(data)
				} else {
					return false
				}
			}
			return true
		})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Controller is not initialized"})
	}
}

func getModelInfo(c *gin.Context) {
	if ctrl != nil && len(ctrl.NodeList) > 0 {
		c.JSON(http.StatusOK, gin.H{"model_path": ctrl.NodeList[0].ModelPath,
			"is_generation": ctrl.NodeList[0].IsGeneration})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Controller is not initialized or no nodes available."})
	}
}
