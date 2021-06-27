package session

import "geeorm/geelog"

func (s *Session) Begin() (err error) {
	geelog.Info("transaction begin")
	if s.tx, err = s.db.Begin(); err != nil {
		geelog.Error(err)
	}
	return
}

func (s *Session) Commit() (err error) {
	geelog.Info("transaction commit")
	if err = s.tx.Commit(); err != nil {
		geelog.Info(err)
	}
	return
}

func (s *Session) Rollback() (err error) {
	geelog.Info("transaction rollback")
	if err = s.tx.Rollback(); err != nil {
		geelog.Info(err)
	}
	return
}
