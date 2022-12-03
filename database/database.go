package database

import (
	"database/sql"
	"fmt"
	"web-app/shared"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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

func Open(connectionString string) {
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

func MigrateDb(connectionString string) error {
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

func OpenMock() (sqlmock.Sqlmock, *sql.DB, error) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  "mock_db",
		DriverName:           "postgres",
		Conn:                 mockDb,
		PreferSimpleProtocol: true,
	})

	db := shared.Unwrap(gorm.Open(dialector, &gorm.Config{}))
	dbCtx = &DbContext{
		db: db,
	}

	return mock, mockDb, nil
}
