package proxy

import (
	"fmt"
	"reflect"
)

type InvocationHandler interface {
	Invoke(proxy *Proxy, method *Method, args ...any) (any, error)
}

type Proxy struct {
	target  any
	methods map[string]*Method
	handle  InvocationHandler
}

func NewProxy(target any, h InvocationHandler) *Proxy {
	typ := reflect.TypeOf(target)
	value := reflect.ValueOf(target)
	methods := make(map[string]*Method)
	for i := 0; i < typ.NumMethod(); i++ {
		methods[typ.Method(i).Name] = &Method{name: typ.Method(i).Name, value: value.Method(i)}
	}
	return &Proxy{target: target, methods: methods, handle: h}
}

func (p *Proxy) InvokeMethod(methodName string, args ...any) (any, error) {
	return p.handle.Invoke(p, p.methods[methodName], args...)
}

type Method struct {
	name  string
	value reflect.Value
}

func (m *Method) Invoke(args ...any) (res []any, err error) {
	defer func() {
		if p := recover(); p != nil {
			// TODO: 处理异常
			err = fmt.Errorf("method %s panic: %v", m.value.Type().Name(), p)
		}
	}()
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}
	res = make([]any, 0)
	for i := 0; i < m.value.Type().NumOut(); i++ {
		res = append(res, m.value.Call(in)[i].Interface())
	}
	return
}
