package orm

import (
	"database/sql"
	"geeorm/dialect"
	"geeorm/geelog"
	"geeorm/session"
	_ "github.com/go-sql-driver/mysql"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func NewEngine(driver, dsn string) (e *Engine, err error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		geelog.Error(err)
		return nil, err
	}
	// ping to make sure the database connection is alive.
	if err = db.Ping(); err != nil {
		geelog.Error(err)
		return nil, err
	}
	geelog.Info("Connect database success")

	dial, ok := dialect.GetDialect(driver)
	if !ok {
		geelog.Errorf("dialect %s Not Found", driver)
		return
	}
	return &Engine{db: db, dialect: dial}, nil
}

func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		geelog.Error("Failed to close database")
	}
	geelog.Info("Close database success")
}

func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dialect)
}
