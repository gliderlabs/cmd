package web

import (
	"crypto/tls"
	"fmt"
	"sync"
)

type CertReloader struct {
	certPath string
	keyPath  string
	cert     *tls.Certificate
	sync.Mutex
}

func NewCertReloader(certPath, keyPath string) (*CertReloader, error) {
	cr := &CertReloader{
		certPath: certPath,
		keyPath:  keyPath,
	}
	err := cr.Reload()
	if err != nil {
		return nil, err
	}
	return cr, nil
}

func (cr *CertReloader) Reload() error {
	c, err := tls.LoadX509KeyPair(cr.certPath, cr.keyPath)
	if err != nil {
		return err
	}
	cr.Lock()
	defer cr.Unlock()
	cr.cert = &c
	return nil
}

func (cr *CertReloader) GetCertificate(clientHello *tls.ClientHelloInfo) (*tls.Certificate, error) {
	cr.Lock()
	defer cr.Unlock()
	if cr.cert == nil {
		return nil, fmt.Errorf("no certificate available")
	}
	return cr.cert, nil
}
