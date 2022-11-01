package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"web-app/shared"
)

const connectionString = "host=localhost user=go_web_app password=go_web_app dbname=go_web_app"

type DbContext struct {
	db       *gorm.DB
	Products IProductsRepository
	Users    IUsersRepository
}

func (ctx DbContext) Close() {
}

func Connect() *DbContext {
	db := shared.Unwrap(gorm.Open(postgres.Open(connectionString), &gorm.Config{}))

	return &DbContext{
		db:       db,
		Products: &gormProductsRepository{db: db},
		Users:    &gormUsersRepository{db: db},
	}
}

func MigrateDb() error {
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&Product{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		return err
	}

	return err
}
