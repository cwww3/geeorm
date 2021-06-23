package orm

import (
	"database/sql"
	"geeorm/geelog"
	"geeorm/session"
)

type Engine struct {
	db *sql.DB
}

func NewEngine(driver, dsn string) (*Engine, error) {
	db, err := sql.Open(driver, dsn)
	if err != nil {
		geelog.Error(err)
		return nil, err
	}
	// ping to make sure the database connection is alive.
	if err = db.Ping();err != nil {
		geelog.Error(err)
		return nil, err
	}
	geelog.Info("Connect database success")
	return &Engine{db: db}, nil
}

func (e *Engine) Close() {
	if err := e.db.Close();err != nil {
		geelog.Error("Failed to close database")
	}
	geelog.Info("Close database success")
}

func (e *Engine) NewSession() *session.Session {
	return session.New(e.db)
}
