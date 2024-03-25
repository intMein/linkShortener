package database

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"sync"
)

type Shortener interface {
	Get(key string) (string, error)
	Has(key string) (bool, error)
	Set(key string, url string) error
}

type InMemoryDB struct {
	links map[string]string
	m     sync.Mutex
}

func (db *InMemoryDB) Get(key string) (string, error) {
	db.m.Lock()
	defer db.m.Unlock()
	v, ok := db.links[key]

	if !ok {
		return "", errors.New("key doesn't exist")
	}

	return v, nil
}

func (db *InMemoryDB) Set(key string, url string) error {
	db.m.Lock()
	defer db.m.Unlock()
	db.links[key] = url

	return nil
}

func (db *InMemoryDB) Has(key string) (bool, error) {
	db.m.Lock()
	defer db.m.Unlock()
	_, ok := db.links[key]

	return ok, nil
}

func NewMemoryDB() Shortener {
	db := &InMemoryDB{}
	db.links = make(map[string]string)

	return db
}

type PostgresDB struct {
	conn *sql.DB
}

func createTable(conn *sql.DB) {
	conn.Exec("CREATE TABLE links (id SERIAL PRIMARY KEY, short VARCHAR(32) UNIQUE, link TEXT)")
}

func NewPostgresDB(user string, password string, database string) (*PostgresDB, error) {
	conn, err := sql.Open("postgres",
		fmt.Sprintf("user=%s password=%s host=127.0.0.1 port=5432 database=%s sslmode=disable",
			user,
			password,
			database))

	if err != nil {
		return nil, err
	}

	createTable(conn)

	db := PostgresDB{}
	db.conn = conn

	return &db, nil
}

func (db *PostgresDB) Set(key string, url string) error {
	stmt, err := db.conn.Prepare("INSERT INTO links (short, link) VALUES ($1, $2)")
	defer stmt.Close()

	if err != nil {
		return err
	}

	_, err = stmt.Exec(key, url)

	return err
}

func (db *PostgresDB) Get(key string) (string, error) {
	stmt, err := db.conn.Prepare("SELECT link FROM links WHERE short = $1")
	defer stmt.Close()

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	res, err := stmt.Query(key)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	var link string

	res.Next()

	res.Scan(&link)

	defer res.Close()

	if link == "" {
		return "", errors.New("not found")
	}

	return link, nil
}

func (db *PostgresDB) Has(key string) (bool, error) {
	_, err := db.Get(key)

	if err != nil {
		return false, nil
	}

	return true, nil
}
