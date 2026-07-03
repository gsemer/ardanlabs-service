// Package keystore implements the auth.KeyLookup interface. THis implements
// an in-memory keystore for JWT support.
package keystore

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"
)

// key represents key information.
type key struct {
	privatePEM string
	publicPEM  string
}

// keyStore represents an in-memory store implementation of the
// KeyLookup interface for us with the auth package.
type keyStore struct {
	store map[string]key
}

// New constructs an empty KeyStore readu for use.
func New() *keyStore {
	return &keyStore{
		store: make(map[string]key),
	}
}

func (ks *keyStore) PrivateKey(kid string) (string, error) {
	key, found := ks.store[kid]
	if !found {
		return "", errors.New("kid lookup failed")
	}

	return key.privatePEM, nil
}

func (ks *keyStore) PublicKey(kid string) (string, error) {
	key, found := ks.store[kid]
	if !found {
		return "", errors.New("kid lookup failed")
	}

	return key.publicPEM, nil
}

// LoadRSAKeys loads a set of RSM PEM files rooted inside of a directory. The
// name of each PEM file will be used as the key id.
// Example: ks.LOadRSAKeys(os.DirFS("/zarf/keys/"))
func (ks *keyStore) LoadRSAKeys(fsys fs.FS) error {
	fn := func(fileName string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("walkdir failure: %w", err)
		}

		if dirEntry.IsDir() {
			return nil
		}

		if path.Ext(fileName) != ".pem" {
			return nil
		}

		file, err := fsys.Open(fileName)
		if err != nil {
			return fmt.Errorf("opening key file: %w", err)
		}
		defer file.Close()

		// limit PEM file size to 1 MB. This should be reasonable for almost any
		// PEM file and prevents shenanigans like linking the file to /dev/random
		// or something like that.
		pem, err := io.ReadAll(io.LimitReader(file, 1024*1024))
		if err != nil {
			return fmt.Errorf("reading auth private key: %w", err)
		}

		privatePEM := string(pem)
		publicPEM, err := toPublicPEM(privatePEM)
		if err != nil {
			return fmt.Errorf("converting private PEM to public: %w", err)
		}

		key := key{
			privatePEM: privatePEM,
			publicPEM:  publicPEM,
		}

		ks.store[strings.TrimPrefix(dirEntry.Name(), ".pem")] = key

		return nil
	}

	if err := fs.WalkDir(fsys, ".", fn); err != nil {
		return fmt.Errorf("walking directory: %w", err)
	}

	return nil
}

func toPublicPEM(privatePEM string) (string, error) {
	block, _ := pem.Decode([]byte(privatePEM))
	if block == nil {
		return "", errors.New("invalid key: Key must be a PEM encoded PKCS1 or PKCS8 key")
	}

	var parsedKey any
	parsedKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return "", err
		}
	}

	pk, ok := parsedKey.(*rsa.PrivateKey)
	if !ok {
		return "", errors.New("key is not a valid RSA private key")
	}

	asn1Bytes, err := x509.MarshalPKIXPublicKey(&pk.PublicKey)
	if err != nil {
		return "", fmt.Errorf("marshaling public key: %w", err)
	}

	// Construct a PEM block for the public key.
	publicBlock := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	var buf bytes.Buffer
	if err := pem.Encode(&buf, &publicBlock); err != nil {
		return "", fmt.Errorf("encoding to public PEM: %w", err)
	}

	return buf.String(), nil
}
