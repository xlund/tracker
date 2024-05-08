package repository

import (
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/corbado/corbado-go"
	"github.com/corbado/corbado-go/pkg/generated/api"
	"github.com/corbado/corbado-go/pkg/generated/common"
	"github.com/xlund/tracker/internal/domain"
)

type corbadoAuthenticator struct {
	Config  *corbado.Config
	Corbado *corbado.Impl
}

func NewCorbadoAuthenticator() domain.Authenticator {
	config, err := corbado.NewConfig(
		os.Getenv("CORBADO_PROJECT_ID"),
		os.Getenv("CORBADO_API_SECRET"),
	)
	if err != nil {
		panic(err)
	}

	sdk, err := corbado.NewSDK(config)

	if err != nil {
		panic(err)
	}
	return &corbadoAuthenticator{Config: config, Corbado: sdk}
}

func (a *corbadoAuthenticator) CreateUser(ctx context.Context, u domain.User) (string, error) {
	var corbadoUser = api.UserCreateReq{
		RequestID: &u.ID,
		Name:      u.Username,
		FullName:  &u.Name,
		Email:     &u.Email,
	}
	usr, err := a.Corbado.Users().Create(ctx, corbadoUser)
	if err != nil {
		return "", err
	}

	return usr.Data.UserID, nil
}

func (a *corbadoAuthenticator) StartWebAuthnRegister(ctx context.Context, u domain.User, r *http.Request) (string, error) {
	log.Default().Println("StartWebAuthnRegister")
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)

	_, err := a.CreateUser(ctx, u)

	res, err := a.Corbado.Passkeys().RegisterStart(ctx, api.WebAuthnRegisterStartReq{
		Username: u.Username,
		ClientInfo: common.ClientInfo{
			RemoteAddress: ip,
			UserAgent:     r.UserAgent(),
		},
	})

	if err != nil {
		return "", err
	}

	return res.PublicKeyCredentialCreationOptions, nil
}

func (a *corbadoAuthenticator) CompleteWebAuthRegister(ctx context.Context, r *http.Request) error {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		return err
	}

	log.Default().Println("CompleteWebAuthRegister")

	// Reading the stringified JSON body from the request
	credential, _ := io.ReadAll(r.Body)

	_, err = a.Corbado.Passkeys().RegisterFinish(ctx, api.WebAuthnFinishReq{
		ClientInfo: common.ClientInfo{
			RemoteAddress: ip,
			UserAgent:     r.UserAgent(),
		},

		// The stringified body from the request
		// does not contain enough data.

		// The rawId, response.attestationObject
		// and response.clientDataJSON is not set.
		PublicKeyCredential: string(credential),
	})

	if err != nil {
		return err
	}

	log.Default().Println("CompleteWebAuthRegister: Finished")
	return nil
}

func (a *corbadoAuthenticator) RemoveUser(ctx context.Context, id string) error {
	// TODO
	return nil
}
