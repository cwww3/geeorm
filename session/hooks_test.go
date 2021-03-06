package session

import (
	"fmt"
	"testing"
)

func TestSession_Hook(t *testing.T) {
	s := testRecordInit(t)
	var users []User
	if err := s.Find(&users); err != nil || len(users) != 2 {
		t.Fatal("failed to query all")
	}
	fmt.Println(users)
}

func (u *User) BeforeQuery(s *Session) error {
	fmt.Println("BeforeQuery")
	fmt.Println(s.RefTable().Name)
	return nil
}
