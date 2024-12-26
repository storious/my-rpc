package model

type RpcResponse struct {
	Data      any    `json:"data"`
	DataType  any    `json:"data_type"`
	Message   string `json:"message"`
	Exception error  `json:"exception"`
}
