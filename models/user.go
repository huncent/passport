package models

import (
	"net/http"
	"time"
)

type UserModel struct {
	Id       int64
	Email    string
	Phone    string
	Password string
	Register time.Time
	Stat     int32
}

func (p *UserModel) Add() {

}
