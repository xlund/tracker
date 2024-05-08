package repository

import (
	"context"
	"encoding/json"
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
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}

	log.Default().Printf("StartWebAuthnRegister::user: %+v", u)
	err = u.Validate()

	if err != nil {
		return "", err
	}

	_, err = a.CreateUser(ctx, u)

	if err != nil {
		return "", err
	}

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
	_, _, err := net.SplitHostPort(r.RemoteAddr)

	if err != nil {
		return err
	}

	log.Default().Println("CompleteWebAuthRegister")
	// parse the json data from r.body
	var cred map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&cred)
	log.Println("This is the data", cred)
	if err != nil {
		log.Default().Println(err)
	}
	return nil

	// req := api.WebAuthnFinishReq{
	// 	ClientInfo: common.ClientInfo{
	// 		RemoteAddress: ip,
	// 		UserAgent:     r.UserAgent(),
	// 	},
	// 	PublicKeyCredential: string(credential),
	// }

	// res, err := a.Corbado.Passkeys().RegisterFinish(ctx, req)
	// if err != nil {
	// 	return err
	// }

	// log.Default().Println("CompleteWebAuthRegister::res")
	// log.Default().Println(res)
	// return nil
}

func (a *corbadoAuthenticator) RemoveUser(ctx context.Context, id string) error {
	_, err := a.Corbado.Users().Delete(ctx, id, api.UserDeleteReq{})
	return err
}
