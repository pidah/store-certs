package main

import (
	"bufio"
	"bytes"
	"compress/gzip"
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
		if strings.HasSuffix(f.Name(), "crt") {
			domainName := strings.TrimSuffix(f.Name(), path.Ext(f.Name()))
			//			fmt.Println(domainName)

			// Open cert and key files on disk.
			certFile, err := os.Open(domainName + ".crt")
			if err != nil {
				fmt.Println(err)
			}

			keyFile, err := os.Open(domainName + ".key")
			if err != nil {
				fmt.Println(err)
			}

			// Read cert and key file content.
			certReader := bufio.NewReader(certFile)
			certContent, err := ioutil.ReadAll(certReader)
			if err != nil {
				fmt.Println(err)
			}

			keyReader := bufio.NewReader(keyFile)
			keyContent, err := ioutil.ReadAll(keyReader)
			if err != nil {
				fmt.Println(err)
			}

			var cb bytes.Buffer
			certCompressed := gzip.NewWriter(&cb)
			certCompressed.Name = certFile.Name()
			certCompressed.Write(certContent)
			certCompressed.Close()

			var kb bytes.Buffer
			keyCompressed := gzip.NewWriter(&kb)
			keyCompressed.Name = keyFile.Name()
			keyCompressed.Write(keyContent)
			keyCompressed.Close()

			certName := prefix + env + "/" + certFile.Name()

			inputCert := &secretsmanager.CreateSecretInput{
				Description:  aws.String(certFile.Name()),
				Name:         aws.String(certName),
				SecretBinary: cb.Bytes(),
			}

			result, err := svc.CreateSecret(inputCert)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(result)

			keyName := prefix + env + "/" + keyFile.Name()

			inputKey := &secretsmanager.CreateSecretInput{
				Description:  aws.String(keyFile.Name()),
				Name:         aws.String(keyName),
				SecretBinary: kb.Bytes(),
			}

			result, err = svc.CreateSecret(inputKey)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(result)

		}
	}
}
