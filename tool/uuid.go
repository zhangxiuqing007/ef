package tool

import uuid "github.com/satori/go.uuid"

//NewUUID 创建一个全新的UUID
func NewUUID() string {
	return uuid.NewV4().String()
}
