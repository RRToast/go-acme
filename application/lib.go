package main

import (
        "crypto/rand"
        "crypto/rsa"
        "crypto/tls"
        "flag"
        "net/http"
	"fmt"
        "strings"
	"encoding/json"

        jose "gopkg.in/square/go-jose.v2"
)

type Header struct {
        Algorithm string `json:"alg"`
        Typ string `json:"typ"` 
        KeyID string `json:"kid,omitempty"`
}

func combine(nonce string) {
        privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
        if err != nil {
                panic(err)
        }

	var signerOpts = jose.SignerOptions{}
	signerOpts.WithHeader("Nonce", nonce)
        signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.PS512, Key: privateKey}, &signerOpts)
        if err != nil {
                panic(err)
        }

	payload := map[string]interface{}{"termsOfServiceAgreed":true, "contact": "mailto:cert-admin@example.org"}
	byts, _ := json.Marshal(payload)
	fmt.Println(string(byts)) 
        signer.Options()
        object, err := signer.Sign(byts)
        if err != nil {
                panic(err)
        }

        serialized := object.FullSerialize()
        println("Payload: ", serialized)

        flag.Parse()

	tlsConfig := &tls.Config{}
        tlsConfig.InsecureSkipVerify = true
        tr := &http.Transport{TLSClientConfig: tlsConfig}
        client := &http.Client{Transport: tr}

        req, err := http.NewRequest("POST", "https://localhost:14000/sign-me-up", strings.NewReader(serialized))
        req.Header.Add("Content-Type", "application/jose+json")

        resp, err := client.Do(req)
        if err != nil {
                println(err.Error())
                panic(err)
        }
        defer resp.Body.Close()
        println("HTTP result status: ", resp.Status)
        println("HTTP result body: ", resp.Body)
}
