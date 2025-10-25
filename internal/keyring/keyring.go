package keyring

import (
	"github.com/zalando/go-keyring"
)

var service = "gomailit"

func SaveCredentials(user, password string) error {
	return keyring.Set(service, user, password)
}

func GetCredentials(user string) (string, error) {
	secret, err := keyring.Get(service, user)
	if err != nil {
		return "", err
	}
	return secret, nil
}
