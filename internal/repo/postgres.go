package repo

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rodeorm/shortener/internal/core"
)

type postgresStorage struct {
	DB               *sqlx.DB    // Драйвер подключения к СУБД
	DBName           string      // Имя БД из конфиг.файла
	ConnectionString string      // Строка подключения из конфиг.файла
	deleteQueue      chan string // канал для удаления URL
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
		fmt.Println("Проблема при создании таблиц", err)
		return err
	}
	return nil
}

// InsertUser сохраняет нового пользователя или возвращает уже имеющегося в наличии
func (s postgresStorage) InsertUser(ctx context.Context, Key int) (*core.User, error) {
	if Key == 0 {
		sqlStatement := `
		INSERT INTO Users (Name)
		VALUES ($1)
		RETURNING id`
		var id int
		err := s.DB.QueryRowContext(ctx, sqlStatement, "yandex").Scan(&id)
		if err != nil {
			fmt.Println("Ошибки при вставке в БД", err)
			return nil, err
		}
		return &core.User{Key: id}, nil
	}
	s.DB.QueryRowContext(ctx, "SELECT ID from Users WHERE ID = $1", fmt.Sprint(Key)).Scan(&Key)

	return &core.User{Key: Key}, nil
}

// InsertShortURL принимает оригинальный URL, генерирует для него ключ, сохраняет соответствие оригинального URL и ключа и возвращает ключ (либо возвращает ранее созданный ключ)
func (s postgresStorage) InsertURL(ctx context.Context, URL, baseURL, userKey string) (string, bool, error) {
	if !core.CheckURLValidity(URL) {
		return "", false, fmt.Errorf("невалидный URL: %s", URL)
	}

	var short string

	// Проверяем на то, что ранее пользователь не сокращал URL
	s.DB.QueryRowContext(ctx, "SELECT short from Urls WHERE original = $1", URL).Scan(&short)
	if short != "" {
		return short, true, nil
	}
	// Вставляем новый URL
	shortKey, err := core.ReturnShortKey(5)
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

func (s postgresStorage) SelectOriginalURL(ctx context.Context, shortURL string) (string, bool, bool, error) {
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

// SelectUserURLHistory возвращает перечень соответствий между оригинальным и коротким адресом для конкретного пользователя
func (s postgresStorage) SelectUserURLHistory(ctx context.Context, Key int) (*[]core.UserURLPair, error) {
	urls := make([]core.UserURLPair, 0, 1)
	rows, err := s.DB.QueryContext(ctx, "SELECT original, short, userID FROM Urls WHERE UserID = $1", fmt.Sprint(Key))
	if err != nil {
		return nil, err
	}

	if rows.Err() != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var pair core.UserURLPair

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

func (s postgresStorage) DeleteURLs(ctx context.Context, URL, userKey string) (bool, error) {
	ch := make(chan string)

	urls := core.GetSliceFromString(URL)

	go func() {
		for _, url := range urls {
			ch <- url
		}
		close(ch)
	}()

	for v := range makeDeletePool(ch) {
		go s.deleteURL(ctx, v, userKey)
	}

	return true, nil
}

func (s postgresStorage) deleteURL(ctx context.Context, url, userKey string) (bool, error) {
	tx, err := s.DB.Begin()

	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx, "UPDATE Urls SET isDeleted = true WHERE short = $1 AND userID = $2")
	if err != nil {
		return false, err
	}
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, url, userKey)
	if err != nil {
		return false, err
	}
	err = tx.Commit()
	if err != nil {
		return false, err
	}
	return true, nil
}
