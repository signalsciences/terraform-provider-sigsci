package provider

import (
	"fmt"

	"github.com/signalsciences/go-sigsci"
)

//Config struct for email and password
type Config struct {
	URL       string
	Email     string
	Password  string
	APIToken  string
	FastlyKey string
}

//Client returns a signal science client
func (c *Config) Client() (interface{}, error) {
	if c.URL != "" {
		sigsci.SetAPIUrl(c.URL)
	}

	if c.Email == "" {
		return nil, fmt.Errorf("email cannot be empty")
	}

	if c.APIToken == "" && c.Password == "" {
		return nil, fmt.Errorf("you must provide either password or api_token")
	}

	var (
		client sigsci.Client
		err    error
	)

	if c.APIToken != "" {
		client = sigsci.NewTokenClient(c.Email, c.APIToken)
	} else if c.Password != "" {
		client, err = sigsci.NewClient(c.Email, c.Password)
		if err != nil {
			return nil, err
		}
	}

	if c.FastlyKey != "" {
		client.SetFastlyKey(c.FastlyKey)
	}

	return client, nil
}
