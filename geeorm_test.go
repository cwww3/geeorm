package orm

import (
	"errors"
	"geeorm/session"
	_ "github.com/go-sql-driver/mysql"
	"testing"
)

func OpenDB(t *testing.T) *Engine {
	t.Helper()
	engine, err := NewEngine("mysql", "root:12345678@tcp(127.0.0.1:3306)/geeorm")
	if err != nil {
		t.Fatal("failed to connect", err)
	}
	return engine
}

type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}

func TestEngine_Transaction(t *testing.T) {
	t.Run("rollback", func(t *testing.T) {
		transactionRollback(t)
	})
	t.Run("commit", func(t *testing.T) {
		transactionCommit(t)
	})
}

func transactionRollback(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	_ = s.Model(&User{}).DropTable()
	_ = s.Model(&User{}).CreateTable()
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		_, err = s.Insert(&User{"Tom", 18})
		_, err = s.Insert(&User{"Tom1", 19})
		return nil, errors.New("Error")
	})
	if err == nil {
		t.Fatal("failed to rollback")
	}
	var count int64
	count, err = s.Count()
	if err != nil || count != 0 {
		t.Fatal("failed to rollback")
	}
}

func transactionCommit(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
	s := engine.NewSession()
	_ = s.Model(&User{}).DropTable()
	_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {
		_ = s.Model(&User{}).CreateTable()
		_, err = s.Insert(&User{"Tom", 18})
		return
	})
	u := &User{}
	_ = s.First(u)
	if err != nil || u.Name != "Tom" {
		t.Fatal("failed to commit")
	}
}
