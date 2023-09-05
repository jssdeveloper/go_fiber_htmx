package main

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Id    uint
	Name  string
	Phone string
}
