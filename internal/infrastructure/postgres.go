package infrastructure

import (
	"fmt"

	"github.com/MidnightHelix/MyGram/internal/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormPostgres interface {
	GetConnection() *gorm.DB
}

type gormPostgresImpl struct {
	master *gorm.DB
}

func NewGormPostgres() GormPostgres {
	return &gormPostgresImpl{
		master: connect(),
	}
}

func connect() *gorm.DB {
	host := "127.0.0.1"
	port := "5432"
	user := "midnight"
	password := "midnight"
	dbname := "mygram"

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&model.User{}, &model.SocialMedia{}, &model.Comment{}, &model.Photo{})
	return db
}

func (g *gormPostgresImpl) GetConnection() *gorm.DB {
	return g.master
}