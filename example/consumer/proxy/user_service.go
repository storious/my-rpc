package proxy

import (
	"bytes"
	"fmt"
	"io"
	"myRPC/example/common/model"
	model2 "myRPC/rpc-easy/model"
	serializer2 "myRPC/rpc-easy/serializer"
	"net/http"
)

// UserService 静态代理
type UserService struct {
}

func (u UserService) GetUser(user model.User) model.User {
	var serializer serializer2.Serializer = serializer2.JsonSerializer{}
	rpcReq := model2.RpcRequest{
		ServiceName:    "UserService",
		MethodName:     "GetUser",
		ParameterTypes: []string{"model.User"},
		Args:           []any{user},
	}
	body, _ := serializer.Serialize(rpcReq)
	res, err := http.Post("http://127.0.0.1:8080/rpc", "application/json", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	resBytes, _ := io.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("http status code: %d, message: %v\n", res.StatusCode, string(resBytes)))
	}
	rpcResp := model2.RpcResponse{}
	_ = serializer.Deserialize(resBytes, &rpcResp)
	var data model.User
	dataBytes, _ := serializer.Serialize(rpcResp.Data)
	_ = serializer.Deserialize(dataBytes, &data)
	return data
}
