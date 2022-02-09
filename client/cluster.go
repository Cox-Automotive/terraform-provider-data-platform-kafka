package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Cluster struct {
	ID                           int         `json:"id"`
	Environment                  Environment `json:"environment"`
	ConfluentID                  string      `json:"confluentId"`
	Name                         string      `json:"name"`
	Cloud                        string      `json:"cloud"`
	Region                       string      `json:"region"`
	Availability                 string      `json:"availability"`
	ClusterType                  string      `json:"clusterType"`
	NetworkingType               string      `json:"networkingType"`
	SecretsManagerSecretArn      string      `json:"secretsManagerSecretArn"`
	SecretsManagerCloudSecretArn string      `json:"secretsManagerCloudSecretArn"`
	RestProxyUrl                 string      `json:"restProxyUrl"`
	BootstrapServers             string      `json:"bootstrapServers"`
	DefaultS3BucketName          string      `json:"defaultS3BucketName"`
}

type Clusters struct {
	Items []Cluster `json:"items"`
}

func getClusters(c Client) ([]Cluster, error) {
	url := fmt.Sprintf("%s/clusters", c.getURL())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var clusters Clusters

	err = json.Unmarshal(res, &clusters)
	if err != nil {
		return nil, err
	}

	return clusters.Items, nil
}

func (c *PrivateClient) GetClusters() ([]Cluster, error) {
	return getClusters(c)
}

func (c *PublicClient) GetClusters() ([]Cluster, error) {
	return getClusters(c)
}

func getCluster(c Client, id int) (*Cluster, error) {
	var cluster Cluster

	url := fmt.Sprintf("%s/clusters/%d", c.getURL(), id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &cluster)
	if err != nil {
		return nil, err
	}

	return &cluster, nil
}

func (c *PrivateClient) GetCluster(id int) (*Cluster, error) {
	return getCluster(c, id)
}

func (c *PublicClient) GetCluster(id int) (*Cluster, error) {
	return getCluster(c, id)
}

func getClusterBy(c Client, paramKey string, paramValue string) (*Cluster, error) {
	url := fmt.Sprintf("%s/clusters", c.getURL())
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

	var clusters Clusters

	err = json.Unmarshal(res, &clusters)
	if err != nil {
		return nil, err
	}

	if len(clusters.Items) == 0 {
		return nil, fmt.Errorf("Cluster with %s %s not found", paramKey, paramValue)
	} else {
		return &clusters.Items[0], nil
	}
}

func (c *PrivateClient) GetClusterByName(name string) (*Cluster, error) {
	return getClusterBy(c, "name", name)
}

func (c *PublicClient) GetClusterByName(name string) (*Cluster, error) {
	return getClusterBy(c, "name", name)
}

func (c *PrivateClient) GetClusterByConfluentID(confluentID string) (*Cluster, error) {
	return getClusterBy(c, "confluentId", confluentID)
}

func (c *PublicClient) GetClusterByConfluentID(confluentID string) (*Cluster, error) {
	return getClusterBy(c, "confluentId", confluentID)
}
