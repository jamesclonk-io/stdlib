package web

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/jamesclonk-io/stdlib/env"
)

type BackendClient struct {
	client *http.Client
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

func (bc *BackendClient) Get(url string) (string, error) {
	resp, err := bc.client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if err := checkResponse(resp, http.StatusOK); err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (bc *BackendClient) Post(url, data string) (string, error) {
	resp, err := bc.client.Post(url, "application/json", strings.NewReader(data))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if err := checkResponse(resp, http.StatusCreated); err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (bc *BackendClient) Put(url, data string) (string, error) {
	req, err := http.NewRequest("PUT", url, strings.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := bc.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if err := checkResponse(resp, http.StatusOK); err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (bc *BackendClient) Delete(url string) (string, error) {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return "", err
	}

	resp, err := bc.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if err := checkResponse(resp, http.StatusNoContent); err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (bc *BackendClient) HttpClient() *http.Client {
	return bc.client
}

func (bc *BackendClient) RootCAFile() string {
	return bc.caFile
}

func (bc *BackendClient) Hostname() string {
	return bc.caHost
}

func checkResponse(resp *http.Response, expectedCode int) error {
	if resp.StatusCode == expectedCode {
		return nil
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return errors.New(string(data))
}
