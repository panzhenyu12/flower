package model

type Validator interface {
	Valid() error
}
