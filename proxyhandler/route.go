package proxyhandler

import (
	"crypto/tls"
	"io/ioutil"
	"encoding/pem"
	pxf "golang.org/x/crypto/pkcs12"
	"strings"
)


// Route config.
type Route struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	Service   string `json:"service"`
	CertFile  string `json:"cert-file"`
	Key       string `json:"key"`
	Transform *Transform `json:"transform"`
}


//
func (r *Route) loadClientCert() (tls.Certificate) {
	if strings.HasSuffix(r.CertFile, ".pem") {
		return r.loadPemCert()
	} else {
		return r.loadPxfCert()
	}
}

//
func (r *Route) loadPemCert() (tls.Certificate) {
	cert, err := tls.LoadX509KeyPair(r.CertFile, r.Key)

	if err != nil {
		panic(err)
	}

	return cert
}

//
func (r *Route) loadPxfCert() (tls.Certificate) {
	var err error
	var data []byte
	data, err = ioutil.ReadFile(r.CertFile)
	if err != nil {
		panic(err)
	}

	var blocks []*pem.Block
	blocks, err = pxf.ToPEM(data, r.Key)
	if err != nil {
		panic(err)
	}

	var pemData []byte
	for _, b := range blocks {
		pemData = append(pemData, pem.EncodeToMemory(b)...)
	}

	cert, err := tls.X509KeyPair(pemData, pemData)
	if err != nil {
		panic(err)
	}

	return cert
}

