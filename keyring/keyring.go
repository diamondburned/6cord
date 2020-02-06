package keyring

import (
	"log"

	keyring "github.com/zalando/go-keyring"
)

const (
	// AppName used for keyrings
	AppName = "6cord"
)

func Get() string      { return Store.Get() }
func Set(token string) { Store.Set(token) }

var Store Storer = defaultKeyring{}

type Storer interface {
	Get() string
	Set(token string)
	Delete()
}

type defaultKeyring struct{}

func (defaultKeyring) Get() string {
	k, err := keyring.Get(AppName, "token")
	if err != nil {
		log.Println("Warning: Could not get token:", err)
	}

	return k
}

func (defaultKeyring) Set(token string) {
	if err := keyring.Set(AppName, "token", token); err != nil {
		log.Println("Warning: Token is not stored:", err)
	}
}

func (defaultKeyring) Delete() {
	if err := keyring.Delete(AppName, "token"); err != nil {
		log.Println("Warning: Token is not deleted:", err)
	}
}
