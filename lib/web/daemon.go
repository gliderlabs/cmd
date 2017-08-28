package web

import "github.com/gliderlabs/comlab/pkg/com"

func (c *Component) DaemonInitialize() error {
	if com.GetString("tls_cert") != "" {
		cr, err := NewCertReloader(com.GetString("tls_cert"), com.GetString("tls_key"))
		if err != nil {
			return err
		}
		c.cert = cr
	}
	return nil
}
