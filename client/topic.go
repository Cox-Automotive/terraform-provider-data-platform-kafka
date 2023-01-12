package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"math"
)

type Topics struct {
	Items []Topic `json:"items"`
}

type Topic struct {
	ID         int          `json:"id,omitempty"`
	Cluster    *Cluster     `json:"cluster,omitempty"`
	Name       string       `json:"name,omitempty"`
	Partitions int          `json:"partitionsCount,omitempty"`
	Config     *TopicConfig `json:"config,omitempty"`
}

type TopicConfig struct {
	ID                int    `json:"id,omitempty"`
	ReplicationFactor int    `json:"replicationFactor,omitempty"`
	MaxMessageBytes   int    `json:"maxMessageBytes,omitempty"`
	CleanupPolicy     string `json:"cleanupPolicy,omitempty"`
	RetentionMs       int    `json:"retentionMs,omitempty"`
	RetentionBytes    int    `json:"retentionBytes,omitempty"`
}

type NewTopic struct {
	ClusterID  int          `json:"clusterId"`
	Name       string       `json:"name"`
	Partitions int          `json:"partitionsCount"`
	Config     *TopicConfig `json:"config"`
}

const topicResourcePath string = "/topics"

func getTopics(c Client) ([]Topic, error) {
	url := fmt.Sprintf("%s%s", c.getURL(), topicResourcePath)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var topics Topics

	err = json.Unmarshal(res, &topics)
	if err != nil {
		return nil, err
	}

	return topics.Items, nil
}

func (c *PrivateClient) GetTopics() ([]Topic, error) {
	return getTopics(c)
}

func (c *PublicClient) GetTopics() ([]Topic, error) {
	return getTopics(c)
}

func getTopic(c Client, id int) (*Topic, error) {
	url := fmt.Sprintf("%s%s/%d", c.getURL(), topicResourcePath, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var topic Topic

	err = json.Unmarshal(res, &topic)
	if err != nil {
		return nil, err
	}

	return &topic, nil
}

func (c *PrivateClient) GetTopic(id int) (*Topic, error) {
	return getTopic(c, id)
}

func (c *PublicClient) GetTopic(id int) (*Topic, error) {
	return getTopic(c, id)
}

func getTopicByNameAndClusterID(c Client, name string, clusterID int) (*Topic, error) {
	url := fmt.Sprintf("%s%s", c.getURL(), topicResourcePath)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("name", name)
	q.Add("cluster.id", strconv.Itoa(clusterID))
	req.URL.RawQuery = q.Encode()

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var topics Topics

	err = json.Unmarshal(res, &topics)
	if err != nil {
		return nil, err
	}

	if len(topics.Items) == 0 {
		return nil, fmt.Errorf("cannot find topic with name %s and cluster ID %d", name, clusterID)
	} else {
		return &topics.Items[0], nil
	}

}

func (c *PrivateClient) GetTopicByNameAndClusterID(name string, clusterID int) (*Topic, error) {
	return getTopicByNameAndClusterID(c, name, clusterID)
}

func (c *PublicClient) GetTopicByNameAndClusterID(name string, clusterID int) (*Topic, error) {
	return getTopicByNameAndClusterID(c, name, clusterID)
}

func createTopic(c Client, t *NewTopic) (*Topic, error) {
	j, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s%s", c.getURL(), topicResourcePath)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}

    retryCap := 3
    res, err := c.doRequest(req)
    for retry := 1; retry <= retryCap; retry++{
        if err != nil {
            fmt.Printf("%s\n\n",err)
            fmt.Println("Waiting to retry...")
            time.Sleep(time.Duration(math.Pow(float64(retry), 2.5)) * time.Second)
            fmt.Printf("Starting retry attempt %d of %d\n", retry, retryCap)
            res, err := c.doRequest(req)
            _,_ = res, err
        } else {
            break
        }
    }
    if err != nil {
        return nil, err
    }

	var topic Topic
	err = json.Unmarshal(res, &topic)
	if err != nil {
		return nil, err
	}

	return &topic, nil
}

func (c *PrivateClient) CreateTopic(t *NewTopic) (*Topic, error) {
	return createTopic(c, t)
}

func (c *PublicClient) CreateTopic(t *NewTopic) (*Topic, error) {
	return createTopic(c, t)
}

func updateTopic(c Client, t *Topic) error {
	url := fmt.Sprintf("%s%s/%d", c.getURL(), topicResourcePath, t.ID)
	//WORKAROUND: Kafka Manager doesn't like the id field in PATCH requests.
	t.ID = 0

	j, err := json.Marshal(t)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", url, bytes.NewBuffer(j))
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *PrivateClient) UpdateTopic(t *Topic) error {
	return updateTopic(c, t)
}

func (c *PublicClient) UpdateTopic(t *Topic) error {
	return updateTopic(c, t)
}

func deleteTopic(c Client, id int) error {
	url := fmt.Sprintf("%s%s/%d", c.getURL(), topicResourcePath, id)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *PrivateClient) DeleteTopic(id int) error {
	return deleteTopic(c, id)
}

func (c *PublicClient) DeleteTopic(id int) error {
	return deleteTopic(c, id)
}
