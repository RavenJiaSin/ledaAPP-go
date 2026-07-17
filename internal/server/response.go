package server


type APIResponse struct {

	Code int `json:"code"`

	Msg string `json:"msg,omitempty"`

	Result interface{} `json:"result,omitempty"`

}

type LiveInferResponse struct {

	ResultImage string `json:"result_image"`
}



type LiveInferResultResponse struct {

	LastResult interface{} `json:"last_result"`

	Timestamp interface{} `json:"timestamp"`
}
