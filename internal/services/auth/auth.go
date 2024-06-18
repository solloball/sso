package auth

import (
    "fmt"
    "context"
    "log/slog"
    "time"
    "errors"

    "golang.org/x/crypto/bcrypt"

    "github.com/solloball/sso/internal/domain/models"
    "github.com/solloball/sso/internal/storage"
    "github.com/solloball/sso/internal/lib/logger/sl"
    "github.com/solloball/sso/internal/lib/jwt"
)

type Auth struct {
    log *slog.Logger
    userSaver UserSaver
    userProvider UserProvider
    appProvider AppProvider
    tokenTTL time.Duration
}

type UserSaver interface {
    SaveUser(
        ctx context.Context,
        email string,
        passHash []byte,
    ) (uid int64, err error)
}

type UserProvider interface {
    User(ctx context.Context, email string) (models.User, error)
    IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
    App(ctx context.Context, appID int) (models.App, error)
}

// New returns a new instance of the Auth service.
func New(
    log *slog.Logger,
    userSaver UserSaver,
    userProvider UserProvider,
    appProvider AppProvider,
    tokenTTL time.Duration,
) *Auth {
    return &Auth {
        log: log,
        userSaver: userSaver,
        userProvider: userProvider,
        appProvider: appProvider,
    }
}

var (
    ErrInvalidData = errors.New("invalid data")
)

func (a *Auth) Login(
    ctx context.Context,
    email string,
    password string,
    appID int,
) (token string, err error) {
    const op = "auth.Login"

    log := a.log.With(
        slog.String("op", op),
        slog.String("email", email),
    )

    log.Info("login user")

    user, err := a.userProvider.User(ctx, email)
    if err != nil {
        if errors.Is(err, storage.ErrAppNotFound) {
            log.Warn("user not found", sl.Err(err))

            return "", fmt.Errorf("%s: %w", op, ErrInvalidData)
        }

        log.Error("failed to get user", sl.Err(err))

        return "", fmt.Errorf("%s: %w", op, err)
    }
    
    if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
        log.Error("invalid data", sl.Err(err))

        return "", fmt.Errorf("%s: %w", op, ErrInvalidData)
    }

    app, err := a.appProvider.App(ctx, appID)
    if err != nil {
        return "", fmt.Errorf("%s: %w", op, err)
    }

    log.Info("user logged in successfully")

    tokenStr, err :=  jwt.NewToken(user, app, a.tokenTTL)
    if err != nil {
        log.Error("failed to make token", sl.Err(err))

        return "", fmt.Errorf("%s: %w", op, err)
    }

    return tokenStr, nil
}

func (a *Auth) Register(
    ctx context.Context,
    email string,
    password string,
) (userID int64, err error) {
    const op = "auth.Register"

    log := a.log.With(
        slog.String("op", op),
        slog.String("email", email),
    )

    log.Info("registering user")

    passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        log.Error("failed to generate password hash", sl.Err(err))

        return 0, fmt.Errorf("%s: %w", op, err)
    }

    id, err := a.userSaver.SaveUser(ctx, email, passHash)
    if err != nil {
        if errors.Is(err, storage.ErrUsrExists) {
            log.Warn("user already exists", sl.Err(err))
            return 0, fmt.Errorf("%s: %w", op, ErrInvalidData)
        }
        log.Error("failed to save user", sl.Err(err))

        return 0, fmt.Errorf("%s: %w", op, err)
    }

    log.Info("user is registered")

    return id, nil
}

func (a *Auth) IsAdmin(
    ctx context.Context,
    userID int64,
) (bool, error) {
	const op = "auth.IsAdmin"

	log := a.log.With(
		slog.String("op", op),
		slog.Int64("user_id", userID),
	)

	log.Info("checking if user is admin")

	isAdmin, err := a.userProvider.IsAdmin(ctx, userID)
	if err != nil {
        if errors.Is(err, storage.ErrAppNotFound) {
            log.Warn("user not found", sl.Err(err))
		    return false, fmt.Errorf("%s: %w", op, ErrInvalidData)
        }
		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("checked if user is admin", slog.Bool("is_admin", isAdmin))

	return isAdmin, nil
}
