package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"web-app/shared"
)

const connectionString = "host=localhost user=go_web_app password=go_web_app dbname=go_web_app"

type DbContext struct {
	db       *gorm.DB
	products IProductsRepository
	users    IUsersRepository
}

func (ctx DbContext) Close() {
}

func (ctx DbContext) Products() IProductsRepository {
	if ctx.products != nil {
		return ctx.products
	}

	ctx.products = &gormProductsRepository{db: ctx.db}

	return ctx.products
}

func (ctx DbContext) Users() IUsersRepository {
	if ctx.users != nil {
		return ctx.users
	}

	ctx.users = &gormUsersRepository{db: ctx.db}

	return ctx.users
}

func Connect() *DbContext {
	db := shared.Unwrap(gorm.Open(postgres.Open(connectionString), &gorm.Config{}))

	return &DbContext{
		db: db,
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
