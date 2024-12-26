package model

type RpcRequest struct {
	ServiceName    string   `json:"service_name"`
	MethodName     string   `json:"method_name"`
	ParameterTypes []string `json:"parameter_types"`
	Args           []any    `json:"args"`
}
