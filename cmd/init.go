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

package cmd

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

// TODO (rluckie) Refactor key generation to use a type with methods and relocate

func generatePrivateKey(bits int) (*rsa.PrivateKey, error) {
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

func encodePrivateKey(privateKey *rsa.PrivateKey) []byte {
	privateKeyDER := x509.MarshalPKCS1PrivateKey(privateKey)

	privateKeyBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privateKeyDER,
	}
	encodedPrivateKey := pem.EncodeToMemory(&privateKeyBlock)
	return encodedPrivateKey
}

func generatePublicKey(privatekey *rsa.PublicKey) ([]byte, error) {
	if publicRsaKey, err := ssh.NewPublicKey(privatekey); err != nil {
		return nil, err
	} else {
		return ssh.MarshalAuthorizedKey(publicRsaKey), nil
	}
}

func writeKeyToFile(key []byte, destination string) error {
	if err := ioutil.WriteFile(destination, key, 0600); err != nil {
		return err
	} else {
		return nil
	}
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize KDK",
	Long:  `Initialize KDK: Create/recreate KDK configuration and pull latest image`,
	Run: func(cmd *cobra.Command, args []string) {

		logger := logrus.New().WithField("command", "init")

		keypairName := "id_rsa"
		keypairDir := path.Join(KdkConfigDir, "ssh")

		privateKeyPath := path.Join(keypairDir, keypairName)

		if _, err := os.Stat(keypairDir); os.IsNotExist(err) {
			if err := os.Mkdir(keypairDir, 0700); err != nil {
				logger.WithField("error", err).Fatal("Failed to create ssh key directory")
			}
		}

		if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
			logger.Info("KDK ssh key not found.")
			logger.Info("Generating KDK ssh key pair...")
			privateKey, err := generatePrivateKey(4096)
			if err != nil {
				logger.WithField("error", err).Fatal("Failed to generate ssh private key")
			}

			publicKeyBytes, err := generatePublicKey(&privateKey.PublicKey)
			if err != nil {
				logger.WithField("error", err).Fatal("Failed to generate ssh public key")
			}

			err = writeKeyToFile(encodePrivateKey(privateKey), privateKeyPath)
			if err != nil {
				logger.WithField("error", err).Fatal("Failed to write ssh private key")
			}

			err = writeKeyToFile([]byte(publicKeyBytes), path.Join(fmt.Sprintf("%s.pub", privateKeyPath)))
			if err != nil {
				logger.WithField("error", err).Fatal("Failed to write ssh public key")
			}
			logger.Info("Successfully generated ssh key pair.")

		} else {
			logger.Info("KDK ssh key pair exists.")
		}

	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
