package sqlite

import (
    "fmt"
    "context"
    "database/sql"
    "errors"
    
    "github.com/mattn/go-sqlite3"
    "github.com/solloball/sso/internal/storage"
    "github.com/solloball/sso/internal/domain/models"
)

type Storage struct {
    db *sql.DB
}

func New(storagePath string) (*Storage, error) {
    const op = "storage.sqlite.New"

    db, err := sql.Open("sqlite3", storagePath)
    if err != nil {
        return nil, fmt.Errorf("%s: :%w", op, err) 
    }

    return &Storage{db: db}, nil
}

func (s *Storage) SaveUser(
    ctx context.Context,
    email string,
    passHash []byte,
) (uid int64, err error) {
    const op = "storage.sqlite.SaveUser"

    stmt, err := s.db.Prepare("INSERT INTO users(email, passHash) VALUES (?,?)")
    if err != nil {
        return 0, fmt.Errorf("%s: %w", op, err)
    }

    res, err := stmt.ExecContext(ctx, email, passHash)
    if err != nil {
        var sqliteErr sqlite3.Error

        if errors.As(err, &sqliteErr) &&
            sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
                return 0, fmt.Errorf("%s: %w", op, storage.ErrUsrExists)
        }

        return 0, fmt.Errorf("%s: %w", op, err)
    }

    id, err := res.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("%s: %w", op, err)
    }

    return id, nil
}

func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
    const op = "storage.sqlite3.User"

    stmt, err := s.db.Prepare(`
        SELECT id, email, pass_hash
        FROM users
        WHERE email == ?`)
    if err != nil {
        return models.User{}, fmt.Errorf("%s: %w", op, err)
    }

    row := stmt.QueryRowContext(ctx, email)

    var user models.User
    
    err = row.Scan(&user.ID, &user.Email, &user.PassHash)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
        }

        return models.User{}, fmt.Errorf("%s: %w", op, err)
    }

    return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
    const op = "storage.sqlite3.IsAdmin"

    stmt, err := s.db.Prepare(`
        SELECT is_admin
        FROM users
        WHERE id == ?`)
    if err != nil {
        return false, fmt.Errorf("%s: %w", op, err)
    }

    row := stmt.QueryRowContext(ctx, userID)

    var res bool
    err = row.Scan(&res)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
        }

        return false, fmt.Errorf("%s: %w", op, err)
    }

    return res, nil
}

func (s *Storage) App(ctx context.Context, id int) (models.App, error) {
	const op = "storage.sqlite.App"

	stmt, err := s.db.Prepare(`
        SELECT id, name, secret
        FROM apps
        WHERE id = ?`,
    )
	if err != nil {
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, id)

	var res models.App
	err = row.Scan(&res.ID, &res.Name, &res.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}

		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return res, nil
}
