package rancher

import (
	"fmt"
	"log"

	"github.com/rancher/go-rancher/catalog"
	rancherClient "github.com/rancher/go-rancher/v2"
)

// Config is the configuration parameters for a Rancher API
type Config struct {
	APIURL    string
	AccessKey string
	SecretKey string
}

// GlobalClient creates a Rancher client scoped to the global API
func (c *Config) GlobalClient() (*rancherClient.RancherClient, error) {
	client, err := rancherClient.NewRancherClient(&rancherClient.ClientOpts{
		Url:       c.APIURL,
		AccessKey: c.AccessKey,
		SecretKey: c.SecretKey,
	})
	if err != nil {
		return nil, err
	}

	log.Printf("[INFO] Rancher Client configured for url: %s", c.APIURL)

	return client, nil
}

// EnvironmentClient creates a Rancher client scoped to an Environment's API
func (c *Config) EnvironmentClient(env string) (*rancherClient.RancherClient, error) {

	globalClient, err := c.GlobalClient()
	if err != nil {
		return nil, err
	}

	project, err := globalClient.Project.ById(env)
	if err != nil {
		return nil, err
	}
	projectURL := project.Links["self"]

	log.Printf("[INFO] Rancher Client configured for url: %s/schemas", projectURL)

	return rancherClient.NewRancherClient(&rancherClient.ClientOpts{
		Url:       projectURL + "/schemas",
		AccessKey: c.AccessKey,
		SecretKey: c.SecretKey,
	})
}

// RegistryClient creates a Rancher client scoped to a Registry's API
func (c *Config) RegistryClient(id string) (*rancherClient.RancherClient, error) {
	client, err := c.GlobalClient()
	if err != nil {
		return nil, err
	}
	reg, err := client.Registry.ById(id)
	if err != nil {
		return nil, err
	}
	if reg == nil {
		return nil, fmt.Errorf("Registry ID %v not found. Check your API key permissions.", id)
	}

	return c.EnvironmentClient(reg.AccountId)
}

// CatalogClient creates a Rancher client scoped to a Catalog's API
func (c *Config) CatalogClient() (*catalog.RancherClient, error) {
	return catalog.NewRancherClient(&catalog.ClientOpts{
		Url:       c.APIURL,
		AccessKey: c.AccessKey,
		SecretKey: c.SecretKey,
	})
}
