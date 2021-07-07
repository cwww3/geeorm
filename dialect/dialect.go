package dialect

import "reflect"

type Dialect interface {
	DataTypeOf(typ reflect.Type) string
	TableExistSQL(tableName string) (string, []interface{})
}

var dialectMap map[string]Dialect

func init() {
	dialectMap = make(map[string]Dialect)
}

func RegisterDialect(name string, dialect Dialect) {
	dialectMap[name] = dialect
}

func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectMap[name]
	return
}
