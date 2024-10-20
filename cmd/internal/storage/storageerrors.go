package storage

import (
	"errors"

	errtext "github.com/RenzoFudo/GoCityMarket/cmd/internal/domain/errors"
)

var ErrInvalidLoginData = errors.New(errtext.InvalidLoginDataError)
var ErrUserNotFound = errors.New(errtext.UserNotFoundError)
var ErrProductNotFound = errors.New(errtext.ProductNotFoundError)
var ErrProductsListEmpty = errors.New(errtext.ProductsListEmptyError)
