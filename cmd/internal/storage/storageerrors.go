package storage

import (
	"errors"

	errtext "github.com/RenzoFudo/http_salem_repo/internal/domain/errors"
)

var ErrInvalidAuthData = errors.New(errtext.InvalidAuthDataError)
var ErrUserNotFound = errors.New(errtext.UserNotFoundError)
var ErrBookNotFound = errors.New(errtext.BookNotFoundError)
var ErrBooksListEmpty = errors.New(errtext.BooksListEmptyError)
