package gateways

type Cache interface {
	Get(key string) (string, bool)
	Set(key string, val string)
	GetAll(key string) []string
	Delete(key string)
	Flush()
}
