package storage

type Storage interface {
	Get(key string) (value []byte, err error)
	Set(key string, value []byte) (err error)
}
