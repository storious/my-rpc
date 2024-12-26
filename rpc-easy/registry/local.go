package registry

/**
 * 本地注册服务侧重根据服务名获取方法
 * 注册中心侧重管理注册的服务
 */

import "sync"

// TODO: redis 注册中心
var mp sync.Map

func Register(serviceName string, handler any) {
	mp.Store(serviceName, handler)
}

func GetService(serviceName string) (any, bool) {
	return mp.Load(serviceName)
}

func RemoveService(serviceName string) {
	mp.Delete(serviceName)
}
