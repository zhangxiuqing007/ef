package tool

import "github.com/satori/uuid"

//NewUUID 创建一个全新的UUID
func NewUUID() string {
	return uuid.NewV4().String()
}
