package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

var env = os.Getenv("env")
var prefix = "ingress/"

type KeyPair struct {
	Cert string `json:"cert"`
	Key  string `json:"key"`
}

func main() {

	if env == "" {
		fmt.Println("You must set the `env` environment variable")
		os.Exit(1)
	}

	files, err := ioutil.ReadDir("./")
	if err != nil {
		log.Fatal(err)
	}

	svc := secretsmanager.New(session.New())

	for _, f := range files {
		if strings.HasSuffix(f.Name(), "cert") {
			domainName := strings.TrimSuffix(f.Name(), path.Ext(f.Name()))
//			fmt.Println(domainName)

			// Open cert and key files on disk.
			certFilename, err := os.Open(domainName + ".cert")
			if err != nil {
				fmt.Println(err)
			}

			keyFilename, err := os.Open(domainName + ".key")
			if err != nil {
				fmt.Println(err)
			}

			// Read cert and key file content.
			certReader := bufio.NewReader(certFilename)
			certContent, err := ioutil.ReadAll(certReader)
			if err != nil {
				fmt.Println(err)
			}

			keyReader := bufio.NewReader(keyFilename)
			keyContent, err := ioutil.ReadAll(keyReader)
			if err != nil {
				fmt.Println(err)
			}

			// Cert Encode as base64.
			certEncoded := base64.StdEncoding.EncodeToString(certContent)

			// Key Encode as base64.
			keyEncoded := base64.StdEncoding.EncodeToString(keyContent)

			secret := KeyPair{certEncoded, keyEncoded}
//			fmt.Println(secret.Key)

			// store encoded data in json
			jsonSecret, err := json.Marshal(secret)
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println(string(jsonSecret))

			secretName := prefix + env + "/" + domainName

			input := &secretsmanager.CreateSecretInput{
				Description:  aws.String(domainName),
				Name:         aws.String(secretName),
				SecretString: aws.String(string(jsonSecret)),
			}

			result, err := svc.CreateSecret(input)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(result)

		}
	}
}
