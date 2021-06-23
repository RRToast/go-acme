package main

import (
        "crypto/rand"
        "crypto/rsa"
        "crypto/tls"
        "flag"
	"io"
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

type dummyNonceSource struct{}

func (n dummyNonceSource) Nonce() (string, error) {
        tr := &http.Transport{
                TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }
        client := &http.Client{Transport: tr}

        res, err := client.Head("https://localhost:14000/nonce-plz")
        if err != nil {
                panic(err)
        }
        ua := res.Header.Get("Replay-Nonce")
        return ua, nil
}

func combine() {
        privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
        if err != nil {
                panic(err)
        }

	var signerOpts = jose.SignerOptions{ NonceSource: dummyNonceSource{}}
	signerOpts.WithHeader("jwk", jose.JSONWebKey{Key: privateKey.Public()})
	signerOpts.WithHeader("url", "https://localhost:14000/sign-me-up")
        signer, err := jose.NewSigner(jose.SigningKey{Algorithm: jose.RS256, Key: privateKey}, &signerOpts)
        if err != nil {
                panic(err)
        }

	var testor [1]string
	testor[0] = "mailto:test.test@test.de"
	payload := map[string]interface{}{"termsOfServiceAgreed":true, "contact": testor}
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
	body, err := io.ReadAll(resp.Body)
        println("HTTP result body: ", string(body))
}
