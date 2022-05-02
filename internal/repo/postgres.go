package repo

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/rodeorm/shortener/internal/logic"
)

type postgresStorage struct {
	DB               *sql.DB // Драйвер подключения к СУБД
	DBName           string  // Имя БД из конфиг.файла
	ConnectionString string  // Строка подключения из конфиг.файла
}

func (s postgresStorage) createTables(ctx context.Context) error {
	_, err := s.DB.ExecContext(ctx,
		"DROP TABLE IF EXISTS Urls;"+
			"DROP TABLE IF EXISTS Users;"+
			"CREATE TABLE IF NOT EXISTS  Users"+
			"("+
			"ID SERIAL PRIMARY KEY"+
			", Name VARCHAR(10) NULL"+
			")"+
			"; CREATE TABLE IF NOT EXISTS  Urls"+
			"("+
			"ID SERIAL PRIMARY KEY"+
			" , Original VARCHAR(1000) NOT NULL "+
			", Short VARCHAR(30) NOT NULL"+
			", UserID	INT  REFERENCES Users (ID) NOT NULL"+
			", CorrelationID varchar(100) NULL"+
			");"+
			"CREATE UNIQUE INDEX IF NOT EXISTS url_unique_idx ON Urls (original, UserID) INCLUDE (short);")
	if err != nil {
		fmt.Println("Проблема при создании таблиц")
		return err
	}
	return nil
}

//InsertUser сохраняет нового пользователя или возвращает уже имеющегося в наличии
func (s postgresStorage) InsertUser(Key int) (*User, error) {
	ctx := context.TODO()

	if Key == 0 {
		sqlStatement := `
		INSERT INTO Users (Name)
		VALUES ($1)
		RETURNING id`
		id := 0
		err := s.DB.QueryRowContext(ctx, sqlStatement, "yandex").Scan(&id)
		if err != nil {
			fmt.Println("Ошибки при вставке в БД", err)
			return nil, err
		}
		return &User{Key: id}, nil
	}
	var key int
	s.DB.QueryRowContext(ctx, "SELECT ID from Users WHERE ID = $1", fmt.Sprint(Key)).Scan(&key)

	return &User{Key: key}, nil
}

//InsertShortURL принимает оригинальный URL, генерирует для него ключ, сохраняет соответствие оригинального URL и ключа и возвращает ключ (либо возвращает ранее созданный ключ)
func (s postgresStorage) InsertURL(URL, baseURL, userKey string) (string, error, bool) {
	if !logic.CheckURLValidity(URL) {

		return "", fmt.Errorf("невалидный URL: %s", URL), false
	}

	ctx := context.TODO()

	var short string

	// Проверяем на то, что ранее пользователь не сокращал URL
	//s.DB.QueryRowContext(ctx, "SELECT short from Urls WHERE UserID = $1 AND original = $2", userKey, URL).Scan(&short)
	s.DB.QueryRowContext(ctx, "SELECT short from Urls WHERE original = $1", URL).Scan(&short)
	if short != "" {
		return short, nil, true
	}
	// Вставляем новый URL
	shortKey, err := logic.ReturnShortKey(5)
	if err != nil {
		fmt.Println(err)
		return "", err, false
	}

	_, err = s.DB.ExecContext(ctx, "INSERT INTO Urls (original, short, userID) SELECT $1, $2, $3", URL, fmt.Sprint(shortKey), userKey)

	if err != nil {
		fmt.Println(err)
		return "", err, false
	}
	// a, err := res.LastInsertId()

	return shortKey, nil, false
}

//SelectOriginalURL принимает на вход короткий URL (относительный, без имени домена), извлекает из него ключ и возвращает оригинальный URL из хранилища
func (s postgresStorage) SelectOriginalURL(shortURL string) (string, bool, error) {
	ctx := context.TODO()
	var original string

	// Проверяем на то, что ранее пользователь не сокращал URL
	s.DB.QueryRowContext(ctx, "SELECT original FROM Urls WHERE short = $1", shortURL).Scan(&original)
	if original != "" {
		return original, true, nil
	}
	return "", false, nil
}

//SelectUserURLHistory возвращает перечень соответствий между оригинальным и коротким адресом для конкретного пользователя
func (s postgresStorage) SelectUserURLHistory(Key int) (*[]UserURLPair, error) {
	urls := make([]UserURLPair, 0, 1)
	ctx := context.TODO()
	fmt.Println("Пользователь", Key)
	rows, err := s.DB.QueryContext(ctx, "SELECT original, short, userID FROM Urls WHERE UserID = $1", fmt.Sprint(Key))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var pair UserURLPair
		err = rows.Scan(&pair.Origin, &pair.Short, &Key)
		if err != nil {
			return nil, err
		}
		urls = append(urls, pair)
	}
	if len(urls) == 0 {
		return nil, fmt.Errorf("нет истории")
	}
	return &urls, nil
}

func (s postgresStorage) CloseConnection() {
	s.DB.Close()
}
