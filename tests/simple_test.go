package tests

import (
    "testing"

    "github.com/golang-jwt/jwt"
    "github.com/brianvoe/gofakeit/v7"
    ssov1 "github.com/solloball/contract/gen/go/sso"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/assert"

    "github.com/solloball/sso/tests/suite"
)

const (
    emptyAppID = 0
    appID = 1
    appSecret = "test"

    passLen = 10
)

func TestRegisterLogin(t *testing.T) {
    ctx, st := suite.New(t)

    email := gofakeit.Email()
    pass := randomFakePassword()

    respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
        Email: email,
        Password: pass,
    })
    require.NoError(t, err)
    assert.NotEmpty(t, respReg.GetUserId())

    respLog, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
        Email: email,
        Password: pass,
        AppId: appID,
    })

    require.NoError(t, err)

    token := respLog.GetToken()
    require.NotEmpty(t, token)

    tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
        return []byte(appSecret), nil
    })
    require.NoError(t, err)

    claims, ok := tokenParsed.Claims.(jwt.MapClaims)
    assert.True(t, ok)

    assert.Equal(t, respReg.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))
}


func TestRegisterLoginDuplicatedRegistration(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	pass := randomFakePassword()

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: pass,
	})
	require.NoError(t, err)
	require.NotEmpty(t, respReg.GetUserId())

	respReg, err = st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: pass,
	})
	require.Error(t, err)
	assert.Empty(t, respReg.GetUserId())
	assert.ErrorContains(t, err, "rpc error: code = AlreadyExists desc = already exists")
}

func TestRegister(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		expectedErr string
	}{
		{
			name:        "Register with Empty Password",
			email:       gofakeit.Email(),
			password:    "",
			expectedErr: "rpc error: code = InvalidArgument desc = password is empty",
		},
		{
			name:        "Register with Empty Email",
			email:       "",
			password:    randomFakePassword(),
			expectedErr: "rpc error: code = InvalidArgument desc = email is empty",
		},
		{
			name:        "Register with Both Empty",
			email:       "",
			password:    "",
			expectedErr: "rpc error: code = InvalidArgument desc = email is empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:    tt.email,
				Password: tt.password,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)

		})
	}
}

func TestLogin(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		email       string
		password    string
		appID       int32
		expectedErr string
	}{
		{
			name:        "Login with Empty Password",
			email:       gofakeit.Email(),
			password:    "",
			appID:       appID,
			expectedErr: "rpc error: code = InvalidArgument desc = password is empty",
        
		},
		{
			name:        "Login with Empty Email",
			email:       "",
			password:    randomFakePassword(),
			appID:       appID,
			expectedErr: "rpc error: code = InvalidArgument desc = email is empty",
		},
		{
			name:        "Login with Both Empty Email and Password",
			email:       "",
			password:    "",
			appID:       appID,
			expectedErr: "rpc error: code = InvalidArgument desc = email is empty",
		},
		{
			name:        "Login with Non-Matching Password",
			email:       gofakeit.Email(),
			password:    randomFakePassword(),
			appID:       appID,
			expectedErr: "rpc error: code = Internal desc = internal error",
		},
		{
			name:        "Login without AppID",
			email:       gofakeit.Email(),
			password:    randomFakePassword(),
			appID:       emptyAppID,
			expectedErr: "rpc error: code = InvalidArgument desc = app_id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Email:    gofakeit.Email(),
				Password: randomFakePassword(),
			})
			require.NoError(t, err)

			_, err = st.AuthClient.Login(ctx, &ssov1.LoginRequest{
				Email:    tt.email,
				Password: tt.password,
				AppId:    tt.appID,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func randomFakePassword() string {
    return gofakeit.Password(true, true, true, true, true, passLen)
}
