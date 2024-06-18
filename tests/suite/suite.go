package suite

import (
    "testing"
    "context"
    "net"
    "strconv"

    ssov1 "github.com/solloball/contract/gen/go/sso"
    "google.golang.org/grpc/credentials/insecure"
    "google.golang.org/grpc"

    "github.com/solloball/sso/internal/config"
)

const (
    host = "localhost"
)

type Suit struct {
    T *testing.T
    Cfg *config.Config
    AuthClient ssov1.AuthClient
}

func New(t *testing.T) (context.Context, *Suit) {
    t.Helper()
    t.Parallel()

    cfg := config.MustLoadPath("../config/local_test.yaml")

    ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPC.Timeout)

    t.Cleanup(func() {
        t.Helper()
        cancelCtx()
    })

    cc, err := grpc.DialContext(
        context.Background(),
		grpcAddress(cfg),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
    ) 
	if err != nil {
		t.Fatalf("grpc server connection failed: %v", err)
	}

    return ctx, &Suit {
        T: t,
        Cfg: cfg,
        AuthClient: ssov1.NewAuthClient(cc),
    }
}

func grpcAddress(cfg *config.Config) string {
    return net.JoinHostPort(host, strconv.Itoa(cfg.GRPC.Port))
}
