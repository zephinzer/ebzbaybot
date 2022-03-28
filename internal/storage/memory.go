package storage

func NewMemory() Storage {
	memory := Memory{}
	return &memory
}

type Memory map[string]interface{}

func (m Memory) Get(key string) (interface{}, error) {
	return m[key], nil
}

func (m Memory) Set(key string, value interface{}) error {
	m[key] = value
	return nil
}
