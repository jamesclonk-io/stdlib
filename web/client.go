package web

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"

	"github.com/jamesclonk-io/stdlib/env"
)

type BackendClient struct {
	*http.Client
	caFile string
	caHost string
}

func NewBackendClient() *BackendClient {
	caFile := env.Get("HTTP_CA_FILE", "")
	caHost := env.Get("HTTP_CA_HOST", "")

	if len(caFile) == 0 {
		panic("Cannot create web.BackendClient without root CA!")
	}
	if len(caHost) == 0 {
		panic("Cannot create web.BackendClient without hostname!")
	}

	tlsConfig := &tls.Config{
		RootCAs:    x509.NewCertPool(),
		ServerName: caHost,
	}
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}

	pemData, err := ioutil.ReadFile(caFile)
	if err != nil {
		panic(err)
	}

	ok := tlsConfig.RootCAs.AppendCertsFromPEM(pemData)
	if !ok {
		panic("Could not load PEM data to root CA!")
	}
	return &BackendClient{client, caFile, caHost}
}

func (bc *BackendClient) GET(url string) (string, error) {
	resp, err := bc.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (bc *BackendClient) RootCAFile() string {
	return bc.caFile
}

func (bc *BackendClient) Hostname() string {
	return bc.caHost
}
