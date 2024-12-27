package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"io"
	common "myRPC/example/common/model"
	"myRPC/rpc-easy/model"
	"myRPC/rpc-easy/registry"
	"myRPC/rpc-easy/serializer"
	"reflect"
)

func handle(c echo.Context) error {
	var serialize serializer.Serializer = serializer.JsonSerializer{}
	var req model.RpcRequest
	body, _ := io.ReadAll(c.Request().Body)
	err := serialize.Deserialize(body, &req)
	if err != nil {
		return err
	}
	// 获取请求
	fmt.Println("request: ", c.Request().Method, c.Request().URL.Path)
	// 处理请求
	if body == nil {
		fmt.Println("body is nil")
		return c.JSON(400, "body is nil")
	}
	service, ok := registry.NewRegistry(registry.GetLocalCache()).GetService(req.ServiceName)
	if !ok {
		fmt.Printf("service: %v not found\n", req.ServiceName)
		return c.JSON(400, "service not found")
	}

	fmt.Printf("method name: %s, args: %v\n", req.MethodName, req.Args)

	method := reflect.ValueOf(service).MethodByName(req.MethodName)

	params := make([]reflect.Value, len(req.Args))

	// 处理方法参数
	for i, v := range req.Args {
		tmp, _ := serialize.Serialize(v)
		var val common.User
		err = serialize.Deserialize(tmp, &val)
		if err != nil {
			return err
		}
		params[i] = reflect.ValueOf(val)
	}
	fmt.Println("parameters :", params, "method type:", method.Type().NumIn())

	result := method.Call(params)[0].Interface().(common.User)

	resp := model.RpcResponse{
		Data:     result,
		DataType: reflect.TypeOf(result).Name(),
		Message:  "ok",
	}
	return c.JSON(200, resp)
}
