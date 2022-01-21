package data

import (
	"fmt"

	_ "github.com/lib/pq"
	"xorm.io/xorm"
)

type User struct {
	ID       int64
	Name     string
	Email    string
	Password string `json:"-"`
}

func createDBEngine() (*xorm.Engine, error) {

	connectionInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", "localhost", 5432, "root", "password", "db")

	engine, err := xorm.NewEngine("postgres", connectionInfo)

	if err != nil {
		return nil, err
	}

	if err := engine.Ping(); err != nil {
		return nil, err
	}

	if err := engine.Sync(new(User)); err != nil {
		return nil, err
	}
	return engine, nil
}

func SetupDatabase() *xorm.Engine {
	engine, err := createDBEngine()
	if err != nil {
		panic(err)
	}
	return engine
}
