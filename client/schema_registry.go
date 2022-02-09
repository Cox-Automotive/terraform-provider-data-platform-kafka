package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SchemaRegistry struct {
	ID                      int         `json:"id"`
	Environment             Environment `json:"environment"`
	ConfluentID             string      `json:"confluentId"`
	Endpoint                string      `json:"endpoint"`
	Cloud                   string      `json:"serviceProvider"`
	Region                  string      `json:"serviceProviderRegion"`
	SecretsManagerSecretArn string      `json:"keysSecretsArn"`
}

type SchemaRegistries struct {
	Items []SchemaRegistry `json:"items"`
}

func getSchemaRegistries(c Client) ([]SchemaRegistry, error) {
	url := fmt.Sprintf("%s/schema-registries", c.getURL())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var schemaRegistries SchemaRegistries

	err = json.Unmarshal(res, &schemaRegistries)
	if err != nil {
		return nil, err
	}

	return schemaRegistries.Items, nil
}

func (c *PrivateClient) GetSchemaRegistries() ([]SchemaRegistry, error) {
	return getSchemaRegistries(c)
}

func (c *PublicClient) GetSchemaRegistries() ([]SchemaRegistry, error) {
	return getSchemaRegistries(c)
}

func getSchemaRegistry(c Client, id int) (*SchemaRegistry, error) {
	var schemaRegistry SchemaRegistry

	url := fmt.Sprintf("%s/schema-registries/%d", c.getURL(), id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &schemaRegistry)
	if err != nil {
		return nil, err
	}

	return &schemaRegistry, nil
}

func (c *PrivateClient) GetSchemaRegistry(id int) (*SchemaRegistry, error) {
	return getSchemaRegistry(c, id)
}

func (c *PublicClient) GetSchemaRegistry(id int) (*SchemaRegistry, error) {
	return getSchemaRegistry(c, id)
}

func getSchemaRegistryBy(c Client, paramKey string, paramValue string) (*SchemaRegistry, error) {
	url := fmt.Sprintf("%s/schema-registries", c.getURL())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add(paramKey, paramValue)
	req.URL.RawQuery = q.Encode()

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var schemaRegistries SchemaRegistries

	err = json.Unmarshal(res, &schemaRegistries)
	if err != nil {
		return nil, err
	}

	if len(schemaRegistries.Items) == 0 {
		return nil, fmt.Errorf("SchemaRegistry with %s %s not found", paramKey, paramValue)
	} else {
		return &schemaRegistries.Items[0], nil
	}
}

func (c *PrivateClient) GetSchemaRegistryByConfluentID(confluentID string) (*SchemaRegistry, error) {
	return getSchemaRegistryBy(c, "confluentId", confluentID)
}

func (c *PublicClient) GetSchemaRegistryByConfluentID(confluentID string) (*SchemaRegistry, error) {
	return getSchemaRegistryBy(c, "confluentId", confluentID)
}
