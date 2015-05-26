package web

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/jamesclonk-io/stdlib/env"
)

type BackendClient struct {
	client         *http.Client
	caFile, caHost string
	user, password string
}

func NewBackendClient() *BackendClient {
	caFile := env.MustGet("JCIO_HTTP_CA_FILE")
	caHost := env.MustGet("JCIO_HTTP_CA_HOST")

	// backends want basic auth with user & password
	user := env.MustGet("JCIO_HTTP_AUTH_USER")
	password := env.MustGet("JCIO_HTTP_AUTH_PASSWORD")

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
	return &BackendClient{client, caFile, caHost, user, password}
}

func (bc *BackendClient) Get(url string) (string, error) {
	req, err := bc.newRequest("GET", url, nil)
	if err != nil {
		return "", err
	}
	return bc.doRequest(req, http.StatusOK)
}

func (bc *BackendClient) Post(url, data string) (string, error) {
	req, err := bc.newRequest("POST", url, strings.NewReader(data))
	if err != nil {
		return "", err
	}
	return bc.doRequest(req, http.StatusCreated)
}

func (bc *BackendClient) Put(url, data string) (string, error) {
	req, err := bc.newRequest("PUT", url, strings.NewReader(data))
	if err != nil {
		return "", err
	}
	return bc.doRequest(req, http.StatusOK)
}

func (bc *BackendClient) Delete(url string) (string, error) {
	req, err := bc.newRequest("DELETE", url, nil)
	if err != nil {
		return "", err
	}
	return bc.doRequest(req, http.StatusNoContent)
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

func (bc *BackendClient) doRequest(req *http.Request, expectedCode int) (string, error) {
	resp, err := bc.client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return "", err
	}

	if err := checkResponse(resp, expectedCode); err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (bc *BackendClient) newRequest(method, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.SetBasicAuth(bc.user, bc.password)

	return req, nil
}

func checkResponse(resp *http.Response, expectedCode int) error {
	if resp.StatusCode == expectedCode {
		return nil
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return errors.New(string(body))
}
