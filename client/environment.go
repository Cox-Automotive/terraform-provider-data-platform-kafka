package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Environment struct {
	ID          int    `json:"id"`
	ConfluentID string `json:"confluentId"`
	Name        string `json:"name"`
	Supplier    string `json:"supplier"`
}

type Environments struct {
	Items []Environment `json:"items"`
}

func getEnvironments(c Client) ([]Environment, error) {
	url := fmt.Sprintf("%s/environments", c.getURL())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var environments Environments

	err = json.Unmarshal(res, &environments)
	if err != nil {
		return nil, err
	}

	return environments.Items, nil
}

func (c *PrivateClient) GetEnvironments() ([]Environment, error) {
	return getEnvironments(c)
}

func (c *PublicClient) GetEnvironments() ([]Environment, error) {
	return getEnvironments(c)
}

func getEnvironment(c Client, id int) (*Environment, error) {
	var environment Environment

	url := fmt.Sprintf("%s/environments/%d", c.getURL(), id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &environment)
	if err != nil {
		return nil, err
	}

	return &environment, nil
}

func (c *PrivateClient) GetEnvironment(id int) (*Environment, error) {
	return getEnvironment(c, id)
}

func (c *PublicClient) GetEnvironment(id int) (*Environment, error) {
	return getEnvironment(c, id)
}

func getEnvironmentBy(c Client, paramKey string, paramValue string) (*Environment, error) {
	url := fmt.Sprintf("%s/environments", c.getURL())
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

	var environments Environments

	err = json.Unmarshal(res, &environments)
	if err != nil {
		return nil, err
	}

	if len(environments.Items) == 0 {
		return nil, fmt.Errorf("Environment with %s %s not found", paramKey, paramValue)
	} else {
		return &environments.Items[0], nil
	}
}

func (c *PrivateClient) GetEnvironmentByName(name string) (*Environment, error) {
	return getEnvironmentBy(c, "name", name)
}

func (c *PublicClient) GetEnvironmentByName(name string) (*Environment, error) {
	return getEnvironmentBy(c, "name", name)
}

func (c *PrivateClient) GetEnvironmentByConfluentID(confluentID string) (*Environment, error) {
	return getEnvironmentBy(c, "confluentId", confluentID)
}

func (c *PublicClient) GetEnvironmentByConfluentID(confluentID string) (*Environment, error) {
	return getEnvironmentBy(c, "confluentId", confluentID)
}
