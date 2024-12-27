package proxy

import (
	"bytes"
	"fmt"
	"io"
	"myRPC/example/common/model"
	model2 "myRPC/rpc-easy/model"
	serializer2 "myRPC/rpc-easy/serializer"
	"net/http"
	"reflect"
)

type ServiceProxy struct {
}

func (s ServiceProxy) Invoke(proxy *Proxy, method *Method, args ...any) (any, error) {
	var serializer serializer2.Serializer = serializer2.JsonSerializer{}
	parametersTypes := getParameterTypes(method.value)
	rpcReq := model2.RpcRequest{
		ServiceName:    reflect.TypeOf(proxy.target).Name(),
		MethodName:     method.name,
		ParameterTypes: parametersTypes,
		Args:           args,
	}
	body, _ := serializer.Serialize(rpcReq)
	// TODO: 使用注册中心和服务发现机制代替硬编码
	res, err := http.Post("http://127.0.0.1:8080/rpc", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resBytes, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http status code: %d, message: %v\n", res.StatusCode, string(resBytes))
	}
	rpcResp := model2.RpcResponse{}
	_ = serializer.Deserialize(resBytes, &rpcResp)
	var data model.User
	dataBytes, _ := serializer.Serialize(rpcResp.Data)
	_ = serializer.Deserialize(dataBytes, &data)
	return data, nil
}

func getParameterTypes(method reflect.Value) []string {
	res := make([]string, method.Type().NumIn())
	for i := 1; i < method.Type().NumIn(); i++ {
		res[i] = method.Type().In(i).Name()
	}
	return res
}
