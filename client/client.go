package client

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client interface {
	getURL() string
	getHTTPClient() *http.Client
	doRequest(req *http.Request) ([]byte, error)

	GetEnvironments() ([]Environment, error)
	GetEnvironment(id int) (*Environment, error)
	GetEnvironmentByName(name string) (*Environment, error)
	GetEnvironmentByConfluentID(confluentID string) (*Environment, error)
	GetSchemaRegistries() ([]SchemaRegistry, error)
	GetSchemaRegistry(id int) (*SchemaRegistry, error)
	GetSchemaRegistryByConfluentID(confluentID string) (*SchemaRegistry, error)
	GetClusters() ([]Cluster, error)
	GetCluster(id int) (*Cluster, error)
	GetClusterByName(name string) (*Cluster, error)
	GetClusterByConfluentID(confluentID string) (*Cluster, error)
	GetTopics() ([]Topic, error)
	GetTopic(id int) (*Topic, error)
	GetTopicByNameAndClusterID(name string, clusterID int) (*Topic, error)
	CreateTopic(t *NewTopic) (*Topic, error)
	UpdateTopic(t *Topic) error
	DeleteTopic(id int) error
}

func doRequest(c Client, req *http.Request) ([]byte, error) {

	res, err := c.getHTTPClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if res.StatusCode == http.StatusOK || res.StatusCode == http.StatusCreated || res.StatusCode == http.StatusNoContent {
		return body, err
	} else {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}
}

type PrivateClient struct {
	HttpClient *http.Client
	URL        string
	Key        string
	OktaGroups string
	Supplier   string
	UserID     string
}

func NewPrivateClient(URL string, Key string, OktaGroups string, Supplier string, UserID string) *PrivateClient {
	return &PrivateClient{
		HttpClient: http.DefaultClient,
		URL:        URL,
		Key:        Key,
		OktaGroups: OktaGroups,
		Supplier:   Supplier,
		UserID:     UserID,
	}
}

func (c *PrivateClient) getURL() string {
	return c.URL
}

func (c *PrivateClient) getHTTPClient() *http.Client {
	return c.HttpClient
}

func (c *PrivateClient) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("CAI-Data-Platform-Key", c.Key)
	req.Header.Set("CAI-Data-Platform-Okta-Groups", c.OktaGroups)
	req.Header.Set("CAI-Data-Platform-Supplier", c.Supplier)
	req.Header.Set("CAI-Data-Platform-User-Id", c.UserID)

	return doRequest(c, req)
}

type PublicClient struct {
	HttpClient  *http.Client
	URL         string
	AccessToken string
}

func NewPublicClient(URL string, AccessToken string) *PublicClient {
	return &PublicClient{
		HttpClient:  http.DefaultClient,
		URL:         URL,
		AccessToken: AccessToken,
	}
}

func (c *PublicClient) getURL() string {
	return c.URL
}

func (c *PublicClient) getHTTPClient() *http.Client {
	return c.HttpClient
}

func (c *PublicClient) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Access-Token", c.AccessToken)

	return doRequest(c, req)
}
