package service

import (
	"log"
	"os"

	"github.com/go-webauthn/webauthn/webauthn"
)

var WA *webauthn.WebAuthn

func WebAuthnInit() {
	var err error
	WA, err = webauthn.New(&webauthn.Config{
		RPDisplayName: os.Getenv("APP_NAME"),
		RPID:          os.Getenv("DNS"),
		RPOrigins:     []string{os.Getenv("ORIGIN")},
	})
	if err != nil {
		log.Fatal("Failed to create WebAuthn from config:", err)
	}
}
