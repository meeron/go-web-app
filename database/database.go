package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"web-app/shared"
)

const connectionString = "host=localhost user=go_web_app password=go_web_app dbname=go_web_app"

var dbCtx *DbContext

type DbContext struct {
	db       *gorm.DB
	products IProductsRepository
	users    IUsersRepository
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

func Open() {
	if dbCtx != nil {
		panic(fmt.Errorf("dbContext already created"))
	}

	db := shared.Unwrap(gorm.Open(postgres.Open(connectionString), &gorm.Config{}))
	dbCtx = &DbContext{
		db: db,
	}
}

func DbCtx() *DbContext {
	if dbCtx == nil {
		panic(fmt.Errorf("dbContext not created"))
	}

	return dbCtx
}

func MigrateDb() error {
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return err
	}

	// Close connection after migration is done
	defer func() {
		sqlDb, err := db.DB()
		if err != nil {
			return
		}

		sqlDb.Close()
	}()

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
