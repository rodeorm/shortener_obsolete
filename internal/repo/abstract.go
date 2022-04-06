package repo

type DB interface {
	InsertShortURL(URL string) (string, error)
	SelectOriginalURL(shortURL string) (string, bool, error)
}
