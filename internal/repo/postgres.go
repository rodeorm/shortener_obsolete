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
			", isDeleted BOOLEAN NOT NULL DEFAULT False"+
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
func (s postgresStorage) InsertURL(URL, baseURL, userKey string) (string, bool, error) {
	if !logic.CheckURLValidity(URL) {

		return "", false, fmt.Errorf("невалидный URL: %s", URL)
	}

	ctx := context.TODO()

	var short string

	// Проверяем на то, что ранее пользователь не сокращал URL
	//s.DB.QueryRowContext(ctx, "SELECT short from Urls WHERE UserID = $1 AND original = $2", userKey, URL).Scan(&short)
	s.DB.QueryRowContext(ctx, "SELECT short from Urls WHERE original = $1", URL).Scan(&short)
	if short != "" {
		return short, true, nil
	}
	// Вставляем новый URL
	shortKey, err := logic.ReturnShortKey(5)
	if err != nil {
		fmt.Println(err)
		return "", false, err
	}

	_, err = s.DB.ExecContext(ctx, "INSERT INTO Urls (original, short, userID) SELECT $1, $2, $3", URL, fmt.Sprint(shortKey), userKey)

	if err != nil {
		fmt.Println(err)
		return "", false, err
	}
	// a, err := res.LastInsertId()

	return shortKey, false, nil
}

func (s postgresStorage) SelectOriginalURL(shortURL string) (string, bool, bool, error) {
	ctx := context.TODO()
	var (
		original  string
		isDeleted bool
	)

	err := s.DB.QueryRowContext(ctx, "SELECT original, isDeleted FROM Urls WHERE short = $1", shortURL).Scan(&original, &isDeleted)
	if err != nil {
		return "", false, false, err
	}
	if isDeleted {
		return original, true, true, nil
	}

	if original != "" {
		return original, true, false, nil
	}

	return "", false, false, nil
}

//SelectUserURLHistory возвращает перечень соответствий между оригинальным и коротким адресом для конкретного пользователя
func (s postgresStorage) SelectUserURLHistory(Key int) (*[]UserURLPair, error) {
	urls := make([]UserURLPair, 0, 1)
	ctx := context.TODO()
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

func (s postgresStorage) DeleteURLs(URL, userKey string) (bool, error) {

	ctx := context.TODO()
	tx, err := s.DB.Begin()
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	//stmt, err := tx.PrepareContext(ctx, "UPDATE Urls SET isDeleted = true WHERE short = $1") // Проще, конечно, через ANY
	stmt, err := tx.PrepareContext(ctx, "UPDATE Urls SET isDeleted = true WHERE short = $1 AND userid = $2") // Проще, конечно, через ANY
	if err != nil {
		return false, err
	}
	defer stmt.Close()
	urls := logic.GetSliceFromString(URL)

	for _, v := range urls {
		_, err = stmt.ExecContext(ctx, v)
		if err != nil {
			return false, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return false, err
	}

	return true, nil
}
