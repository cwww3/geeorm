package main

import (
	"fmt"
	orm "geeorm"
	"geeorm/geelog"
)

func main() {
	e, _ := orm.NewEngine("mysql", "root:12345678@tcp(localhost:3306)/geeorm")
	defer e.Close()
	s := e.NewSession()
	_, _ = s.Raw("DROP TABLE IF EXISTS user;").Exec()
	_, _ = s.Raw("create table user(name varchar(255));").Exec()
	_, _ = s.Raw("create table user(name varchar(255));").Exec()
	_, _ = s.Raw("insert into user values (?), (?)", "Tom", "Sam").Exec()

	s.Raw("select * from user where name = ?;", "Tom")
	rows,err  := s.QueryRows()
	if err != nil {
		geelog.Error(err)
	}
	defer rows.Close()
	for rows.Next() {
		var u string
		err = rows.Scan(&u)
		fmt.Println(u)
	}
}
