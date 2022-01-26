package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Dyleme/image-coverter/internal/model"
)

type jwtToken struct {
	Token string `json:"jwt"`
}

const fileFlag = 0o755

func credentialsFromFile(filePath string) ([]byte, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("credentaials from file: %w", err)
	}

	reqBody, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("credentaials from file: %w", err)
	}

	return reqBody, nil
}

func credentialsFromArgs(nickname, password, email string) ([]byte, error) {
	fmt.Printf("reqeust with %s and %s\n", nickname, password)

	user := model.User{
		Nickname: nickname,
		Password: password,
		Email:    email,
	}

	js, err := json.Marshal(user)
	if err != nil {
		return nil, fmt.Errorf("credentaials from args: %w", err)
	}

	return js, nil
}

const pathToJWT = "/../../.token/.jwt"

func saveJWT(b []byte) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("save jwt: %w", err)
	}

	file, err := os.OpenFile(dir+pathToJWT, os.O_TRUNC|os.O_CREATE|os.O_RDWR, fileFlag)
	if err != nil {
		return fmt.Errorf("save jwt: %w", err)
	}

	_, err = file.Write(b)
	if err != nil {
		return fmt.Errorf("save jwt: %w", err)
	}

	return nil
}

func deleteJWT() error {
	return os.Remove(".jwt")
}

const (
	AuthorizationHeader = "Authorization"

	BearerToken = "Bearer"
)

func auth(r *http.Request) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("save jwt: %w", err)
	}

	file, err := os.OpenFile(dir+pathToJWT, 0, os.ModeType)
	if err != nil {
		return fmt.Errorf("auth: %w", err)
	}

	b, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("auth: %w", err)
	}

	token, err := getToken(b)
	if err != nil {
		return fmt.Errorf("auth: %w", err)
	}

	r.Header.Add(AuthorizationHeader, BearerToken+" "+token)

	return nil
}

func getToken(b []byte) (string, error) {
	jwt := jwtToken{}

	if err := json.Unmarshal(b, &jwt); err != nil {
		return "", fmt.Errorf("get token: %w", err)
	}

	if jwt.Token == "" {
		return "", fmt.Errorf("invalid password")
	}

	return jwt.Token, nil
}
