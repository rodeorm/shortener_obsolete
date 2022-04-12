package repo

import "github.com/golobby/container"

type AbstractStorage interface {
	InsertShortURL(URL string) (string, error)
	SelectOriginalURL(shortURL string) (string, bool, error)
}

func NewStorage(filePath string) AbstractStorage {
	container.Singleton(func() AbstractStorage {

		if filePath != "" {
			storage := fileStorage{filePath: filePath}
			storage.CheckFile()
			return &storage
		}
		ots := make(map[string]string)
		sto := make(map[string]string)
		storage := memoryStorage{originalToShort: ots, shortToOriginal: sto}
		return &storage

	})
	var st AbstractStorage
	container.Make(&st)
	return st
}
