package storage

type Storage interface {
	Get(key string) (value interface{}, err error)
	Set(key string, value interface{}) (err error)
}
