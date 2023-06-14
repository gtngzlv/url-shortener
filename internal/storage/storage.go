package storage

var Storage = make(map[string]string)

func GetFromStorage(key string) string {
	return Storage[key]
}
