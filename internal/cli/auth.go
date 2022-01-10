package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/Dyleme/image-coverter/internal/model"
)

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

func saveJWT(b []byte) error {
	file, err := os.OpenFile(".jwt", os.O_TRUNC|os.O_CREATE, os.ModeTemporary)
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
	file, err := os.OpenFile(".jwt", 0, os.ModeType)
	if err != nil {
		return fmt.Errorf("auth: %w", err)
	}

	b, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("auth: %w", err)
	}

	jwt := struct {
		Token string `json:"jwt"`
	}{}

	err = json.Unmarshal(b, &jwt)
	if err != nil {
		return fmt.Errorf("auth: %w", err)
	}

	r.Header.Add(AuthorizationHeader, BearerToken+" "+jwt.Token)

	return nil
}
