/*
 * Copyright 2020 PeerGum
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package layer

import (
	"bufio"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"io/ioutil"
	"log"
	"os"
)

var pubKey *rsa.PublicKey
var privKey *rsa.PrivateKey

type Key interface{}

func genKey() (key *rsa.PrivateKey, err error) {
	if key, err = rsa.GenerateKey(rand.Reader, 1024); err == nil {
		//if pemBlock,err := pem x509.EncryptPEMBlock(rand.Reader,"RSA PRIVATE KEY",x509.MarshalPKCS1PrivateKey(key),[]byte{},x509.PEMCipherDES); err == nil {
		if keyFile, err := os.Create(*keyPath); err == nil {
			writer := bufio.NewWriter(keyFile)
			pem.Encode(writer, &pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: x509.MarshalPKCS1PrivateKey(key),
			})
			writer.Flush()
			keyFile.Close()
		}
		if keyFile, err := os.Create(*keyPath + ".pub"); err == nil {
			writer := bufio.NewWriter(keyFile)
			pem.Encode(writer, &pem.Block{
				Type:  "RSA PUBLIC KEY",
				Bytes: x509.MarshalPKCS1PublicKey(&key.PublicKey),
			})
			writer.Flush()
			keyFile.Close()
		}
		//}
	}
	return key, err
}

func getPrivateKey(privKeyPath string) (*rsa.PrivateKey, error) {
	if privKeyPath == "" {
		return genKey()
	}

	priv, err := ioutil.ReadFile(privKeyPath)
	if err != nil {
		return genKey()
	}

	privPem, _ := pem.Decode(priv)
	var privPemBytes []byte
	if privPem != nil && privPem.Type != "RSA PRIVATE KEY" {
		log.Fatalf("Private key is not an RSA private key. Pem Type [%s]", privPem.Type)
	}

	if *keyChallenge != "" {
		privPemBytes, err = x509.DecryptPEMBlock(privPem, []byte(*keyChallenge))
	} else {
		privPemBytes = privPem.Bytes
	}

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(privPemBytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(privPemBytes); err != nil { // note this returns type `interface{}`
			log.Print("Can't parse RSA private key", err)
			return genKey()
		}
	}

	var key *rsa.PrivateKey
	var ok bool

	if key, ok = parsedKey.(*rsa.PrivateKey); !ok {
		return genKey()
	}
	return key, nil
}

func (content Content) encrypt(key Key, label []byte) (Content, error) {
	rng := rand.Reader
	encryptedContent, err := rsa.EncryptOAEP(crypto.SHA256.New(), rng, key.(*rsa.PublicKey), content, label)
	if err != nil {
		return nil, err
	}
	return encryptedContent, nil
}

func (content Content) sign() ([]byte, error) {
	rng := rand.Reader
	hash := sha256.Sum256(content)
	signedContent, err := rsa.SignPKCS1v15(rng, privKey, crypto.SHA256, hash[:])
	if err != nil {
		return nil, err
	}
	return signedContent, nil
}
