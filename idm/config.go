package idm

import (
	"crypto/tls"
	"log"
	"net/http"

	sdk "github.com/tehwalris/go-freeipa/freeipa"
)

//Config for config
type Config struct {
	Host     string
	Username string
	Password string
	Insecure bool
}

//NewClient creates a Redhat IDM client scoped to the global API
func (c *Config) NewClient() (*sdk.Client, error) {
	tspt := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: c.Insecure,
		},
	}

	client, err := sdk.Connect(c.Host, tspt, c.Username, c.Password)
	if err != nil {
		return nil, err
	}

	log.Printf("[INFO] Redhat IDM Client configured for host: %s", c.Host)

	return client, nil
}
