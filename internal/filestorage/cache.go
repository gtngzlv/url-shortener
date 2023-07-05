package filestorage

var CacheStorage = make(map[string]string)

func WriteToCache(fullURL, shortURL string) {
	CacheStorage[shortURL] = fullURL
}

func GetFromCache(shortURL string) string {
	return CacheStorage[shortURL]
}
