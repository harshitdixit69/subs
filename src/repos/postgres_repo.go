package repos

import (
	"fmt"
	"os"
	"sync"

	"github.com/fliptable-io/subscription-service/src/utils"
	"github.com/fliptable-io/subscription-service/src/utils/errors"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresRepo struct {
	*sqlx.DB
}

func safeguardParams(val interface{}) interface{} {
	if val == nil {
		return make(map[string]interface{})
	}
	return val
}

type querer interface {
	PrepareNamed(query string) (*sqlx.NamedStmt, error)
}

func queryDb(sql querer, dest interface{}, query string, params interface{}) (err error) {
	statement, err := sql.PrepareNamed(query)
	if err != nil {
		return errors.UnknownError.Consume(err)
	}

	defer func() {
		cerr := statement.Close()
		if cerr != nil && err == nil {
			err = cerr
		}
	}()

	params = safeguardParams(params)

	if dest == nil {
		_, err = statement.Exec(params)
		return errors.UnknownError.Consume(err)
	}

	if utils.IsSlicePointer(dest) {
		err = statement.Select(dest, params)
		return errors.UnknownError.Consume(err)
	}

	err = statement.Get(dest, params)
	if err != nil {
		return errors.NotFoundError.Consume(err)
	}
	return nil
}

func (p *PostgresRepo) Query(dest interface{}, query string, params interface{}) error {
	return queryDb(p, dest, query, params)
}

func (p *PostgresRepo) QueryT(tx *sqlx.Tx, dest interface{}, query string, params interface{}) error {
	err := queryDb(tx, dest, query, params)
	return err
}

var lock = &sync.Mutex{}
var postgresRepo *PostgresRepo

func NewPostgresRepo() *PostgresRepo {
	if postgresRepo != nil {
		return postgresRepo
	}

	lock.Lock()
	defer lock.Unlock()

	if postgresRepo != nil {
		return postgresRepo
	}

	config := fmt.Sprintf(
		"host=%v port=%v dbname=%v user=%v password=%v sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
	)

	db, err := sqlx.Connect("postgres", config)
	if err != nil {
		_ = errors.UnknownError.Consume(err)
	}

	postgresRepo = &PostgresRepo{db}

	return postgresRepo
}
