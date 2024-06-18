package storage

import "errors"

var(
    ErrUsrExists = errors.New("user already exist")
    ErrNotFound= errors.New("user not found")
    ErrAppNotFound = errors.New("app not found")
)

func (*Storage) SaveUser(
    ctx context.Context,
    email string,
    passHash []byte,
) (uid int64, err error) {
}
