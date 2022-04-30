package repo

type AbstractStorage interface {
	InsertURL(URL, baseURL, userKey string) (string, error)  // Сохраняет соответствие между оригинальным и коротким адресом
	SelectOriginalURL(shortURL string) (string, bool, error) // Возвращает оригинальный адрес на основании короткого
	InsertUser(Key int) (*User, error)                       // Ищет существующего пользователя и возвращает его, если находит, если не находит то, сохраняет нового и возвращает его
	SelectUserURLHistory(Key int) (*[]UserURLPair, error)    // Возвращает перечень соответствий между оригинальным и коротким адресом для конкретного пользователя
}

// NewStorage определяет место для хранения данных
func NewStorage(filePath string) AbstractStorage {
	if filePath != "" {
		return initFileStorage(filePath)
	}
	storage := initMemoryStorage()
	return storage
}

func initMemoryStorage() *memoryStorage {
	ots := make(map[string]string)
	sto := make(map[string]string)
	usr := make(map[int]*User)
	usrURL := make(map[int]*[]UserURLPair)
	storage := memoryStorage{originalToShort: ots, shortToOriginal: sto, users: usr, userURLPairs: usrURL}
	return &storage
}

func initFileStorage(filePath string) *fileStorage {
	storage := fileStorage{filePath: filePath}
	storage.CheckFile(filePath)
	return &storage
}
