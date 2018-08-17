// Copyright © 2018 Cisco Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
)

func GeneratePrivateKey(bits int) (*rsa.PrivateKey, error) {
	if privateKey, err := rsa.GenerateKey(rand.Reader, bits); err != nil {
		return nil, err
	} else {
		if err := privateKey.Validate(); err != nil {
			return nil, err
		} else {
			return privateKey, nil
		}
	}
}

func EncodePrivateKey(privateKey *rsa.PrivateKey) []byte {
	privateKeyDER := x509.MarshalPKCS1PrivateKey(privateKey)

	privateKeyBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privateKeyDER,
	}
	encodedPrivateKey := pem.EncodeToMemory(&privateKeyBlock)
	return encodedPrivateKey
}

func GeneratePublicKey(privatekey *rsa.PublicKey) ([]byte, error) {
	if publicRsaKey, err := ssh.NewPublicKey(privatekey); err != nil {
		return nil, err
	} else {
		return ssh.MarshalAuthorizedKey(publicRsaKey), nil
	}
}

func WriteKeyToFile(key []byte, destination string) error {
	if err := ioutil.WriteFile(destination, key, 0600); err != nil {
		return err
	} else {
		return nil
	}
}
