package controllers

import (
	"bytes"
	"encoding/gob"
	"testing"
)

func TestSessionGob(t *testing.T) {
	s := &Session{
		Sid:          "asdasdasdasdasdasd",
		PostSortType: 0,
		UserID:       2,
		UserName:     "小王八",
	}
	buffer := new(bytes.Buffer)
	if err := gob.NewEncoder(buffer).Encode(s); err != nil {
		t.Fatal(err)
	}
	t.Log(string(buffer.Bytes()))
	newS := new(Session)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(newS); err != nil {
		t.Fatal(err)
	}
	t.Log(newS)
}
