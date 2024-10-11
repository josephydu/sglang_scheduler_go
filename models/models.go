package models

type NodeInfo struct {
	Ip                 string `json:"ip"`
	Port               int    `json:"port"`
	ModelPath          string `json:"model_path"`
	IsGeneration       bool   `json:"is_generation"`
	ControllerInfoPort int    `json:"controller_info_port"`
}

type Request interface {
	ToMap() map[string]interface{}
}

type CompletionRequest struct {
	Model       string  `json:"model"`
	Prompt      string  `json:"prompt"`
	Temperature float64 `json:"temperature"`
	BestOf      float64 `json:"best_of"`
	MaxTokens   int     `json:"max_tokens"`
	Stream      bool    `json:"stream"`
	IgnoreEos   bool    `json:"ignore_eos"`
}

func (r CompletionRequest) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"model":       r.Model,
		"prompt":      r.Prompt,
		"temperature": r.Temperature,
		"best_of":     r.BestOf,
		"max_tokens":  r.MaxTokens,
		"stream":      r.Stream,
		"ignore_eos":  r.IgnoreEos,
	}
}

type SamplingParams struct {
	SkipSpecialTokens          bool     `json:"skip_special_tokens"`
	SpacesBetweenSpecialTokens bool     `json:"spaces_between_special_tokens"`
	MaxNewTokens               int      `json:"max_new_tokens"`
	MinNewTokens               int      `json:"min_new_tokens"`
	Stop                       []string `json:"stop"`
	StopTokenIds               []int    `json:"stop_token_ids"`
	Temperature                float64  `json:"temperature"`
	TopP                       float64  `json:"top_p"`
	TopK                       int      `json:"top_k"`
	MinP                       float64  `json:"min_p"`
	FrequencyPenalty           float64  `json:"frequency_penalty"`
	PresencePenalty            float64  `json:"presence_penalty"`
	IgnoreEos                  bool     `json:"ignore_eos"`
	Regex                      *string  `json:"regex"`
	JsonSchema                 *string  `json:"json_schema"`
}

type GenerateRequest struct {
	Text           string         `json:"text"`
	SamplingParams SamplingParams `json:"sampling_params"`
}

func (r GenerateRequest) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"text": r.Text,
		"sampling_params": map[string]interface{}{
			"skip_special_tokens":           r.SamplingParams.SkipSpecialTokens,
			"spaces_between_special_tokens": r.SamplingParams.SpacesBetweenSpecialTokens,
			"max_new_tokens":                r.SamplingParams.MaxNewTokens,
			"min_new_tokens":                r.SamplingParams.MinNewTokens,
			"stop":                          r.SamplingParams.Stop,
			"stop_token_ids":                r.SamplingParams.StopTokenIds,
			"temperature":                   r.SamplingParams.Temperature,
			"top_p":                         r.SamplingParams.TopP,
			"top_k":                         r.SamplingParams.TopK,
			"min_p":                         r.SamplingParams.MinP,
			"frequency_penalty":             r.SamplingParams.FrequencyPenalty,
			"presence_penalty":              r.SamplingParams.PresencePenalty,
			"ignore_eos":                    r.SamplingParams.IgnoreEos,
			"regex":                         r.SamplingParams.Regex,
			"json_schema":                   r.SamplingParams.JsonSchema,
		},
	}
}
