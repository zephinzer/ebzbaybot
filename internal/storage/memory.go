package storage

func NewMemory() Storage {
	memory := Memory{}
	return &memory
}

type Memory map[string][]byte

func (m Memory) Get(key string) ([]byte, error) {
	return m[key], nil
}

func (m Memory) Set(key string, value []byte) error {
	m[key] = value
	return nil
}
