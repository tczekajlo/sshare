package ssh

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
	"log"

	"golang.org/x/crypto/ssh"
)

// Keys stores generated keys
type Keys struct {
	PublicKeyPEM  string
	PublicKey     rsa.PublicKey
	PrivateKeyPEM string
	PrivateKey    *rsa.PrivateKey
}

// Init creates a new pair of RSA keys
func (s *Keys) Init() {
	reader := rand.Reader
	bitSize := 4096

	key, err := rsa.GenerateKey(reader, bitSize)
	printError(err)
	s.PrivateKey = key

	s.PublicKey = key.PublicKey

	s.savePEMKey()
	s.savePublicPEMKey()
}

func (s *Keys) savePEMKey() {
	s.PrivateKeyPEM = string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(s.PrivateKey),
	}))

}

func (s *Keys) savePublicPEMKey() {
	asn1Bytes, err := asn1.Marshal(s.PublicKey)
	printError(err)

	s.PublicKeyPEM = string(pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}))
}

// GetPublicKey returns an OpenSSH authorized_key
func (s *Keys) GetPublicKey() (string, error) {
	publicRsaKey, err := ssh.NewPublicKey(&s.PrivateKey.PublicKey)
	if err != nil {
		return "", err
	}

	return string(ssh.MarshalAuthorizedKey(publicRsaKey)), nil
}

func printError(err error) {
	if err != nil {
		log.Fatalln("Fatal error ", err.Error())
	}
}
