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

type TxFunc func(*session.Session) (interface{}, error)

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

func (e *Engine) Transaction(tx TxFunc) (result interface{}, err error) {
	s := e.NewSession()
	if err := s.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p) // re-throw panic after Rollback
		} else if err != nil {
			_ = s.Rollback() // err is non-nil; don't change it
		} else {
			defer func() {
				if err != nil {
					_ = s.Rollback()
				}
			}()
			err = s.Commit() // err is nil; if Commit returns error update err
		}
	}()
	return tx(s)
}
