package models

import (
	"github.com/go-webauthn/webauthn/webauthn"
)

func (u *User) WebAuthnID() []byte {
	return []byte(u.ID.Hex())
}
func (u *User) WebAuthnName() string {
	return u.Mobile
}
func (u *User) WebAuthnDisplayName() string {
	return u.Name + " " + u.Family
}
func (u *User) WebAuthnIcon() string {
	return u.Image
}
func (u *User) WebAuthnCredentials() []webauthn.Credential {
	creds := []webauthn.Credential{}
	for _, c := range u.FingerTokens {
		creds = append(creds, webauthn.Credential{
			ID:        c.CredentialID,
			PublicKey: c.PublicKey,
			Authenticator: webauthn.Authenticator{
				SignCount: c.SignCount,
			},
		})
	}
	return creds
}
