package registry

/**
 * 本地注册服务侧重根据服务名获取方法
 * 注册中心侧重管理注册的服务
 */

type Registry interface {
	Register(string, any)
	GetService(string) (any, bool)
	RemoveService(string)
}

type registry struct {
	cache Cache
}

func NewRegistry(cache Cache) Registry {
	return &registry{
		cache: cache,
	}
}

func (r *registry) Register(serviceName string, handler any) {
	r.cache.SetItem(serviceName, handler)
}

func (r *registry) GetService(serviceName string) (any, bool) {
	service := r.cache.GetItem(serviceName)
	if service == nil {
		return nil, false
	}
	return service, true
}

func (r *registry) RemoveService(serviceName string) {
	r.cache.SetItem(serviceName, nil)
}
