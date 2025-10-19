package utils

// MapHelper map辅助工具
type MapHelper map[string]interface{}

func (m MapHelper) GetString(key string, defaultVal ...string) string {
	def := ""
	if len(defaultVal) > 0 {
		def = defaultVal[0]
	}
	if val, ok := m[key].(string); ok {
		return val
	}
	return def
}

func (m MapHelper) GetBool(key string, defaultVal ...bool) bool {
	def := false
	if len(defaultVal) > 0 {
		def = defaultVal[0]
	}
	if val, ok := m[key].(bool); ok {
		return val
	}
	return def
}

func (m MapHelper) GetInt(key string, defaultVal ...int) int {
	def := 0
	if len(defaultVal) > 0 {
		def = defaultVal[0]
	}
	if val, ok := m[key].(int); ok {
		return val
	}
	// 兼容 float64（JSON解析常见）
	if val, ok := m[key].(float64); ok {
		return int(val)
	}
	return def
}

func (m MapHelper) Get(key string) interface{} {
	return m[key]
}
