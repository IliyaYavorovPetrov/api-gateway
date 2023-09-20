package gateways

type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, val interface{})
	GetAll(key string) []interface{}
	Delete(key string)
	Flush()
}
