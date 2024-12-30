package test

import (
	"errors"
	"myRPC/rpc-core/model"
	"myRPC/rpc-core/serialization"
	"testing"
)

func TestMyJson(t *testing.T) {
	var data = model.RpcRequest{
		ServiceName:    "HelloService",
		MethodName:     "Hello",
		ParameterTypes: []int{model.TypeString},
		Args:           []any{"world"},
	}
	var serializer = serialization.JsonSerializer{}
	serialize, err := serializer.Serialize(data)
	if err != nil {
		t.Error(err)
		return
	}
	var res model.RpcRequest
	err = serializer.Deserialize(serialize, &res)
	if err != nil {
		t.Error(err)
		return
	}
	if err = equal(data, res); err != nil {
		t.Error(err)
	}
}

func equal(lhs model.RpcRequest, rhs model.RpcRequest) error {
	if lhs.ServiceName != rhs.ServiceName {
		return errors.New("ServiceName not equal")
	}

	if lhs.MethodName != rhs.MethodName {
		return errors.New("MethodName not equal")
	}

	if len(lhs.ParameterTypes) != len(rhs.ParameterTypes) {
		return errors.New("ParameterTypes length not equal")
	}

	for i := 0; i < len(lhs.ParameterTypes); i++ {
		if lhs.ParameterTypes[i] != rhs.ParameterTypes[i] {
			return errors.New("ParameterTypes not equal")
		}
	}
	if len(lhs.Args) != len(rhs.Args) {
		return errors.New("args length not equal")
	}
	for i := 0; i < len(lhs.Args); i++ {
		if lhs.Args[i] != rhs.Args[i] {
			return errors.New("args not equal")
		}
	}
	return nil
}
