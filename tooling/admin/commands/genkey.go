package commands

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

// GenKey creates a new RSA key pair
func GenKey(keysFolder string) error {
	// Generate unique folder for key
	kid := uuid.New()
	directoryPath := filepath.Join(keysFolder, kid.String())
	if err := os.MkdirAll(directoryPath, 0755); err != nil {
		return fmt.Errorf("creating key directory: %w", err)
	}

	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return fmt.Errorf("generating private key: %w", err)
	}

	privateBlock := pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	privateFilePath := filepath.Join(directoryPath, "private.pem")
	privateFile, err := os.Create(privateFilePath)
	if err != nil {
		return fmt.Errorf("creating private key file: %w", err)
	}
	defer privateFile.Close()

	if err := pem.Encode(privateFile, &privateBlock); err != nil {
		return fmt.Errorf("encoding private key to file: %w", err)
	}

	// Generate public key
	asn1Bytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return fmt.Errorf("marshalling public key: %w", err)
	}

	publicBlock := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	publicFilePath := filepath.Join(directoryPath, "public.pem")
	publicFile, err := os.Create(publicFilePath)
	if err != nil {
		return fmt.Errorf("creating public key file: %w", err)
	}

	if err := pem.Encode(publicFile, &publicBlock); err != nil {
		return fmt.Errorf("encoding public key to file: %w", err)
	}

	fmt.Println("Generated RSA key pair: kid = ", kid)

	return nil
}
