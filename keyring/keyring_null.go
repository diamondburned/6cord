// +build nokeyring

package keyring

import (
	"log"
)

func init() {
	Store = nullKeyring{}
}

type nullKeyring struct{}

func (nullKeyring) Get() string {
	log.Println("Warning: 6cord compiled without keyring")
	return ""
}

func (nullKeyring) Set(token string) {
	log.Println("Warning: 6cord compiled without keyring")
}

func (nullKeyring) Delete() {
	log.Println("Warning: 6cord compiled without keyring")
}
