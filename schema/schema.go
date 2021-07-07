package schema

import (
	"geeorm/dialect"
	"go/ast"
	"reflect"
)

// Field represent a column
type Field struct {
	Name string
	Type string
	Tag  string
}

// Schema represent a table
type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Field
	FieldNames []string
	fieldMap   map[string]*Field
}

func (s *Schema) GetField(name string) *Field {
	return s.fieldMap[name]
}

func (s *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, name := range s.FieldNames {
		fieldValues = append(fieldValues, destValue.FieldByName(name).Interface())
	}
	return fieldValues
}

func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(),
		fieldMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		//name := []rune(p.Name)
		//if !p.Anonymous && len(name) > 0 && name[0] >= 'A' && name[0] <= 'Z'{
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				// New() reflect.Type -> reflect.Value   通过反射生成对象 与new()类似 生成指针指向零值对象
				// https://colobu.com/2019/01/29/go-reflect-performance/
				// https://vimsky.com/examples/usage/reflect-new-function-in-golang-with-examples.html
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			if v, ok := p.Tag.Lookup("geeorm"); ok {
				field.Tag = v
			}
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}
