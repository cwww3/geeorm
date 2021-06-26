package session

import (
	"errors"
	"geeorm/clause"
	"reflect"
)

func (s *Session) Insert(values ...interface{}) (int64, error) {
	var recordValues []interface{}
	for _, value := range values {
		table := s.Model(value).refTable
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		recordValues = append(recordValues, table.RecordValues(value))
	}
	// 打散
	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	// 打散
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) Find(value interface{}) error {
	// hooks BeforeQuery
	if err := s.doBeforeQuery(); err != nil {
		return err
	}

	destSlice := reflect.Indirect(reflect.ValueOf(value))
	// 获取切片的单个元素的类型
	// destSlice.Type() -> []session.User
	// destSlice.Elem() -> session.User
	destType := destSlice.Type().Elem()
	table := s.Model(reflect.New(destType).Elem().Interface()).RefTable()

	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).QueryRows()
	if err != nil {
		return err
	}
	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var values []interface{}
		for _, name := range table.FieldNames {
			// dest的每一个字段的指针
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		// 解析dest的每一个字段
		if err := rows.Scan(values...); err != nil {
			return err
		}
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}

// Update kv support map[string]interface{}
// also support kv list: "Name", "Tom", "Age", 18, ...
func (s *Session) Update(kv ...interface{}) (int64, error) {
	m, ok := kv[0].(map[string]interface{})
	if !ok {
		m = make(map[string]interface{})
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
	}
	s.clause.Set(clause.UPDATE, s.RefTable().Name, m)
	sql, values := s.clause.Build(clause.UPDATE, clause.WHERE)
	result, err := s.Raw(sql, values...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) Delete() (int64, error) {
	s.clause.Set(clause.DELETE, s.RefTable().Name)
	sql, values := s.clause.Build(clause.DELETE, clause.WHERE)
	result, err := s.Raw(sql, values...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.RefTable().Name)
	sql, values := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, values...).QueryRow()
	var count int64
	if err := row.Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (s *Session) Where(desc interface{}, args ...interface{}) *Session {
	s.clause.Set(clause.WHERE, append(append([]interface{}{}, desc), args...)...)
	return s
}

func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.ORDERBY, desc)
	return s
}

func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

func (s *Session) First(value interface{}) error {
	dest := reflect.ValueOf(value).Elem()
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("NOT FOUND")
	}
	dest.Set(destSlice.Index(0))
	return nil
}
