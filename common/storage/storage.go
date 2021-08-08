package storage

// Storage is the interface implemented by objects that can store data with string keys
type Storage interface {
	Get(string) interface{}
	Store(string, interface{})
	Delete(string)
}

// RuntimeStorage returns runtime storage.
func RuntimeStorage() Storage {
	return make(rtStore)
}

type rtStore map[string]interface{}

func (r rtStore) Get(key string) interface{} {
	return r[key]
}
func (r rtStore) Store(key string, val interface{}) {
	r[key] = val
}
func (r rtStore) Delete(key string) {
	if _, ok := r[key]; ok {
		delete(r, key)
	}
}
