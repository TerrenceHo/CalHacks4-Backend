package controllers

type PublicError interface {
	error
	Public() string
}
