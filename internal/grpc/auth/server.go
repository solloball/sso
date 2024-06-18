package auth

import (
    "context"
    "errors"

    "google.golang.org/grpc"
    "google.golang.org/grpc/status"
    "google.golang.org/grpc/codes"

    ssov1 "github.com/solloball/contract/gen/go/sso"
    "github.com/solloball/sso/internal/services/auth"
)

type Auth interface {
    Login(
        ctx context.Context,
        email string,
        password string,
        appID int,
    ) (token string, err error)
    Register(
        ctx context.Context,
        email string,
        password string,
    ) (userID int64, err error)
    IsAdmin(ctx context.Context, userID int64) (res bool, err error)
}

type serverAPI struct {
    ssov1.UnimplementedAuthServer 
    auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
    ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

const (
    emptyValue = 0
)

func (s *serverAPI) Login(
    ctx context.Context,
    req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
    if err := validateDataLogin(req); err != nil {
        return nil, err
    }

    token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()));
    if err != nil {
        if errors.Is(err, auth.ErrInvalidData) {
            return nil, status.Error(codes.InvalidArgument, "invalid argument")
        }
        return nil, status.Error(codes.Internal, "internal error")
    }

    return &ssov1.LoginResponse{
        Token: token,
    }, nil
}

func (s *serverAPI) Register(
    ctx context.Context,
    req *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
    if err := validateDataRegister(req); err != nil {
        return nil, err
    }

    userID, err := s.auth.Register(ctx, req.GetEmail(), req.GetPassword())
    if err != nil {
        if errors.Is(err, auth.ErrInvalidData) {
            return nil, status.Error(codes.AlreadyExists, "already exists")
        }
        return nil, status.Error(codes.Internal, "internal error")
    }

    return &ssov1.RegisterResponse{
        UserId: userID,
    }, nil
}

func (s *serverAPI) IsAdmin(
    ctx context.Context,
    req *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponse, error) {
    if err := validateDataIsAdmin(req); err != nil {
        return nil, err
    }

    res, err := s.auth.IsAdmin(ctx, req.GetUserId())
    if err != nil {
        if errors.Is(err, auth.ErrInvalidData) {
            return nil, status.Error(codes.NotFound, "not found")
        }
        return nil, status.Error(codes.Internal, "internal error")
    }

    return &ssov1.IsAdminResponse{
       IsAdmin: res, 
    }, nil
}

func validateDataLogin(req *ssov1.LoginRequest) error {
    if req.GetEmail() == "" {
        return status.Error(codes.InvalidArgument, "email is empty")
    }

    if req.GetPassword() == "" {
        return status.Error(codes.InvalidArgument, "password is empty")
    }

    if req.GetAppId() == emptyValue {
        return status.Error(codes.InvalidArgument, "app_id is required")
    }

    return nil
}

func validateDataRegister(req *ssov1.RegisterRequest) error {
    if req.GetEmail() == "" {
        return status.Error(codes.InvalidArgument, "email is empty")
    }

    if req.GetPassword() == "" {
        return status.Error(codes.InvalidArgument, "password is empty")
    }

    return nil
}

const (
    emptyUserIDValue = 0
)

func validateDataIsAdmin(req *ssov1.IsAdminRequest) error {
    if req.GetUserId() == emptyUserIDValue {
        return status.Error(codes.InvalidArgument, "email is empty")
    }

    return nil
}

