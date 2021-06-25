package clause

import "strings"

type TYPE int

const (
	INSERT = iota
	SELECT
	ORDERBY
	WHERE
	LIMIT
)

type Clause struct {
	sql     map[TYPE]string
	sqlVars map[TYPE][]interface{}
}

func (c *Clause) Set(name TYPE, values ...interface{}) {
	if c.sql == nil {
		c.sql = make(map[TYPE]string)
		c.sqlVars = make(map[TYPE][]interface{})
	}
	// notice values...
	sql, vars := generators[name](values...)
	c.sql[name] = sql
	c.sqlVars[name] = vars
}

func (c *Clause) Build(orders ...TYPE) (string, []interface{}){
	var sqls []string
	var vars []interface{}
	for _, order := range orders {
		if sql,ok := c.sql[order]; ok {
			sqls = append(sqls, sql)
			vars = append(vars, c.sqlVars[order]...)
		}
	}
	return strings.Join(sqls," "), vars
}
