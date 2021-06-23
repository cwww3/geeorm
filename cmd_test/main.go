package main

import (
	"fmt"
	orm "geeorm"
	"geeorm/geelog"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	e, _ := orm.NewEngine("mysql", "root:123456@tcp(127.0.0.1:3307)/geeorm")
	defer e.Close()
	s := e.NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS user;").Exec()
	_, _ = s.Raw("create table user(name varchar(255));").Exec()
	_, _ = s.Raw("create table user(name varchar(255));").Exec()
	_, _ = s.Raw("insert into user values (?), (?)", "Tom", "Sam").Exec()

	s.Raw("select * from user where name = ?;","Tom")
	row := s.QueryRow()
	var u string
	err := row.Scan(&u)
	if err != nil {
		geelog.Error(err)
	}
	fmt.Println(u)
}
