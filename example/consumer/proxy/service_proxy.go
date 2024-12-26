package proxy

import (
	"bytes"
	"io"
	model2 "myRPC/rpc-easy/model"
	serializer2 "myRPC/rpc-easy/serializer"
	"net/http"
	"reflect"
)

type ServiceProxy struct {
}

func (s ServiceProxy) Invoke(_ *Proxy, method *Method, args ...any) (any, error) {
	var serializer serializer2.Serializer = serializer2.JsonSerializer{}
	parametersTypes := getParameterTypes(method.value)
	rpcReq := model2.RpcRequest{
		ServiceName:    method.value.Type().In(0).Name(),
		MethodName:     method.value.Type().Name(),
		ParameterTypes: parametersTypes,
		Args:           args,
	}
	body, _ := serializer.Serialize(rpcReq)
	// TODO: 使用注册中心和服务发现机制代替硬编码
	res, _ := http.Post("localhost:8080/rpc", "application/json", bytes.NewBuffer(body))

	resBytes, _ := io.ReadAll(res.Body)
	rpcResp := model2.RpcResponse{}
	_ = serializer.Deserialize(resBytes, &rpcResp)
	return rpcResp.Data, nil
}

func getParameterTypes(method reflect.Value) []string {
	res := make([]string, method.Type().NumIn())
	for i := 1; i < method.Type().NumIn(); i++ {
		res[i] = method.Type().In(i).Name()
	}
	return res
}
