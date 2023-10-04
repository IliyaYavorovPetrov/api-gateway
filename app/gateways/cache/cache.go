package cache

var instances map[string]interface{}

func init() {
	instances = make(map[string]interface{})
}

func Pool() *map[string]interface{} {
	return &instances
}
