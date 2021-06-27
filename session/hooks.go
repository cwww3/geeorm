package session

import "reflect"

type BeforeQuery interface {
	BeforeQuery(s *Session) error
}

func (s *Session) doBeforeQuery() error {
	_, ok := s.RefTable().Model.(BeforeQuery)
	if ok {
		method := reflect.ValueOf(s.RefTable().Model).MethodByName("BeforeQuery")
		result := method.Call([]reflect.Value{reflect.ValueOf(s)})
		err := result[0].Interface()
		if err != nil {
			return err.(error)
		}
	}
	return nil
}

