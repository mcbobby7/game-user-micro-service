package httperror

type HTTPError struct {
	Message  string      `json:"message"`
	HasError string      `json:"hasError"`
	Data     interface{} `json:"data"`
}