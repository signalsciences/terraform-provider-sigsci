package provider

import (
	"fmt"
	"github.com/signalsciences/go-sigsci"
)

//Config struct for email and password
type Config struct {
	Email    string
	Password string
	APIToken string
}

//Client returns a signal science client
func (c *Config) Client() (interface{}, error) {
	if c.Email == "" {
		return nil, fmt.Errorf("Email cannot be empty", c.Email)
	}
	if c.APIToken != "" {
		return sigsci.NewTokenClient(c.Email, c.APIToken), nil
	} else if c.Password != "" {
		return sigsci.NewClient(c.Email, c.Password)
	}
	return nil, fmt.Errorf("You must provide either password or api_token")
}
