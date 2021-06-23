package main

import (
        "crypto/tls"
        "encoding/base64"
        "encoding/json"
        "net/http"
        "os"
)

type Message struct {
        alg   string
        jwk   string
        nonce string
        url   string
}

func main() {
        nonce := getNonce()
        combine(nonce)
        m := Message{"ES256", "todo", nonce, "https://localhost:14000/sign-me-up"}
        b, err := json.Marshal(m)
        if err != nil {
                panic(err)
        }
        encoder := base64.NewEncoder(base64.StdEncoding, os.Stdout)
        encoder.Write(b)
        encoder.Close()
}

func getNonce() string {
        tr := &http.Transport{
                TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        }
        client := &http.Client{Transport: tr}

        res, err := client.Head("https://localhost:14000/nonce-plz")
        if err != nil {
                panic(err)
        }
        ua := res.Header.Get("Replay-Nonce")
        return ua
}

