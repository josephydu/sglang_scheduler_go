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
	Model       string  `json:"model,omitempty"`
	Prompt      string  `json:"prompt,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
	BestOf      float64 `json:"best_of,omitempty"`
	MaxTokens   int     `json:"max_tokens,omitempty"`
	Stream      bool    `json:"stream,omitempty"`
	IgnoreEos   bool    `json:"ignore_eos,omitempty"`
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
	SkipSpecialTokens          bool     `json:"skip_special_tokens,omitempty"`
	SpacesBetweenSpecialTokens bool     `json:"spaces_between_special_tokens,omitempty"`
	MaxNewTokens               int      `json:"max_new_tokens,omitempty"`
	MinNewTokens               int      `json:"min_new_tokens,omitempty"`
	Stop                       []string `json:"stop,omitempty"`
	StopTokenIds               []int    `json:"stop_token_ids,omitempty"`
	Temperature                float64  `json:"temperature,omitempty"`
	TopP                       float64  `json:"top_p,omitempty"`
	TopK                       int      `json:"top_k,omitempty"`
	MinP                       float64  `json:"min_p,omitempty"`
	FrequencyPenalty           float64  `json:"frequency_penalty,omitempty"`
	PresencePenalty            float64  `json:"presence_penalty,omitempty"`
	IgnoreEos                  bool     `json:"ignore_eos,omitempty"`
	Regex                      *string  `json:"regex,omitempty"`
	JsonSchema                 *string  `json:"json_schema,omitempty"`
}

type GenerateRequest struct {
	Text           string         `json:"text,omitempty"`
	SamplingParams SamplingParams `json:"sampling_params,omitempty"`
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
