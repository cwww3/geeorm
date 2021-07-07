package clause

import (
	"fmt"
	"strings"
)

const (
	INSERT = iota
	SELECT
	ORDERBY
	WHERE
	LIMIT
	VALUES
	UPDATE
	DELETE
	COUNT
)

type TYPE int

type generator func(values ...interface{}) (string, []interface{})

var generators = make(map[TYPE]generator)

func init() {
	generators[INSERT] = _insert
	generators[SELECT] = _select
	generators[ORDERBY] = _orderBy
	generators[WHERE] = _where
	generators[LIMIT] = _limit
	generators[VALUES] = _values
	generators[UPDATE] = _update
	generators[DELETE] = _delete
	generators[COUNT] = _count
}

// 两个参数 第一个参数为表明表名  第二个参数为map[string]interface{}类型存储更新的字段
func _update(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	m := values[1].(map[string]interface{})
	var keys []string
	var vars []interface{}
	for k, v := range m {
		keys = append(keys, k+" = ?")
		vars = append(vars, v)
	}
	return fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(keys, ", ")), vars
}

// 一个参数 表名
func _delete(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	return fmt.Sprintf("DELETE FROM %s", tableName), []interface{}{}
}

// 一个参数 表名
func _count(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	return _select(tableName, []string{"COUNT(*)"})
}

// 两个参数 第一个表名 第二个[]interface类型  表示字段名
func _insert(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %s (%v) ", tableName, fields), []interface{}{}
}

// 可变参数  每个参数为[]interface类型 表示一条记录的字段值
func _values(values ...interface{}) (string, []interface{}) {
	var sql strings.Builder
	var bindStr string
	var vars []interface{}
	sql.WriteString("VALUES ")
	// 多条记录
	for i, value := range values {
		v := value.([]interface{})
		if len(bindStr) == 0 {
			bindStr = getBindVars(len(v))
		}
		sql.WriteString(fmt.Sprintf("(%v)", bindStr))
		if i+1 != len(values) {
			sql.WriteString(", ")
		}
		vars = append(vars, v...)
	}
	return sql.String(), vars
}

func getBindVars(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ",")
}

// 两个参数 第一个表名 第二个[]string类型 表示字段名
func _select(values ...interface{}) (string, []interface{}) {
	fmt.Println(len(values))
	tableName := values[0].(string)
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("SELECT %v FROM %s", fields, tableName), []interface{}{}
}

// 一个参数
func _limit(values ...interface{}) (string, []interface{}) {
	// 只取第一个 避免多传出错
	return "LIMIT ?", []interface{}{values[0]}
}

// 一个参数
func _orderBy(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("ORDER BY %s", values[0]), []interface{}{}
}

// 多个参数 第一个表示条件语句 其他的表示条件值
func _where(values ...interface{}) (string, []interface{}) {
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %s", desc), vars
}
