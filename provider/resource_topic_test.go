package provider

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"

	"coxautoinc.com/data-platform/kafka-manager/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func init() {
	resource.AddTestSweepers("kafkamanager_topic", &resource.Sweeper{
		Name: "kafkamanager_topic",
		F:    testSweepTopics,
		// Dependencies: []string{""},
	})
}

const topicNamePrefix = "tf_acc_test_"

func testSweepTopics(region string) error {
	// NOTE: region is not used at this moment and is only needed to conform to Terraform testing API.
	c, err := sharedPrivateClientForRegion(region)
	if err != nil {
		return err
	}

	topics, err := c.GetTopics()
	if err != nil {
		return err
	}

	for _, t := range topics {
		if strings.HasPrefix(t.Name, topicNamePrefix) {
			log.Printf("Deleting Topic %s", t.Name)

			if err := c.DeleteTopic(t.ID); err != nil {
				return err
			}
		}
	}

	return nil
}

func TestUnmarshalNewTopic(t *testing.T) {
	resourcedatatest := schema.TestResourceDataRaw(t, resourceTopic().Schema, nil)
	resourcedatatest.Set("cluster_id", "10")
	resourcedatatest.Set("name", "test-topic")

	goodTestData := &client.NewTopic{
		ClusterID:  10,
		Name:       "test-topic",
		Partitions: 0,
		Config:     &client.TopicConfig{},
	}

	unmarshaltopic, err := unmarshalNewTopic(resourcedatatest)

	if err != nil {
		t.Fatalf("error expanding perms: %v", err)
	}
	compare := reflect.DeepEqual(goodTestData, unmarshaltopic)
	if !compare {
		t.Fatalf("Error matching, expected: %#v and got %#v", goodTestData, unmarshaltopic)
	}

}

func TestMarshalTopic(t *testing.T) {
	environment := &client.Environment{
		ID:          10,
		ConfluentID: "env-xxxx",
		Name:        "new-env",
		Supplier:    "abc-123",
	}
	cluster := &client.Cluster{
		ID:                           10,
		Environment:                  *environment,
		ConfluentID:                  "lkc-wrjpg",
		Name:                         "cluster1",
		Cloud:                        "aws",
		Region:                       "az1",
		Availability:                 "LOW",
		ClusterType:                  "Standard",
		NetworkingType:               "BASIC",
		SecretsManagerSecretArn:      "ARN:XXXXXX",
		SecretsManagerCloudSecretArn: "ARN:XXXXX",
		RestProxyUrl:                 "http://rest.com",
		BootstrapServers:             "https://1.1.1.1",
		DefaultS3BucketName:          "Bucket",
	}
	topicConfig := &client.TopicConfig{
		ID:                1,
		ReplicationFactor: 3,
		MaxMessageBytes:   16,
		CleanupPolicy:     "delete",
		RetentionMs:       30,
		RetentionBytes:    16,
	}
	topic := &client.Topic{
		ID:         10,
		Name:       "test-topic",
		Partitions: 6,
		Cluster:    cluster,
		Config:     topicConfig,
	}
	goodTestData := map[string]interface{}{
		"id":                 "10",
		"name":               "test-topic",
		"cluster_id":         "10",
		"partitions":         6,
		"max_message_bytes":  16,
		"cleanup_policy":     "delete",
		"retention_ms":       30,
		"retention_bytes":    16,
		"replication_factor": 3,
	}
	result, err := marshalTopic(topic)

	if err != nil {
		t.Fatalf("error expanding perms: %v", err)
	}
	compare := reflect.DeepEqual(goodTestData, result)
	if !compare {
		t.Fatalf("Error matching, expected: %#v and got %#v", goodTestData, result)
	}
}

func TestAccTopicCreation(t *testing.T) {
	// TODO: Replace hardcoded values with something.
	awsEnvironment := "cinp"
	supplierCode := "dev-2"

	topic := &client.Topic{
		Name: fmt.Sprintf("test_acc_%s", acctest.RandString(10)),
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKafkaManagerTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(
					testAccCheckKafkaManagerTopicConfig_basic,
					awsEnvironment,
					supplierCode,
					topic.Name,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKafkaManagerTopicExists("kafkamanager_topic.foobar", topic),
					resource.TestCheckResourceAttr("kafkamanager_topic.foobar", "name", topic.Name),
					resource.TestCheckResourceAttr("kafkamanager_topic.foobar", "partitions", "6"),
					resource.TestCheckResourceAttr("kafkamanager_topic.foobar", "replication_factor", "3"),
				),
			},
		},
	})
}

const testAccCheckKafkaManagerTopicConfig_basic = `
data "kafkamanager_cluster" "bazqux" {
	name = "data-platform-%s-%s-cluster"
}

resource "kafkamanager_topic" "foobar" {
	name               = "%s"
	cluster_id         = data.kafkamanager_cluster.bazqux.id
	partitions         = 6
	replication_factor = 3
	retention_ms       = 30
	}`

func testAccCheckKafkaManagerTopicExists(rn string, topic *client.Topic) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("not found: %s", rn)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no topic ID is set")
		}

		c := testAccProvider.Meta().(*client.PrivateClient)

		topic_id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Invalid topic ID: %s", err)
		}

		got, err := c.GetTopic(topic_id)
		if err != nil {
			return err
		}
		if got.Name != topic.Name {
			return fmt.Errorf("wrong topic found, want %q got %q", topic.Name, got.Name)
		}

		// get the computed topic details
		*topic = *got
		return nil
	}
}

func testAccCheckKafkaManagerTopicDestroy(s *terraform.State) error {
	c := testAccProvider.Meta().(*client.PrivateClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "kafkamanager_topic" {
			continue
		}

		topic_id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Invalid topic ID: %s", err)
		}

		_, err = c.GetTopic(topic_id)
		if err == nil {
			return fmt.Errorf("Topic still exists")
		}

		notFoundErr := "topics.does_not_exist"
		expectedErr := regexp.MustCompile(notFoundErr)
		if !expectedErr.Match([]byte(err.Error())) {
			return fmt.Errorf("expected %s, got %s", notFoundErr, err)
		}
	}

	return nil
}

func TestAccTopicModification_UpdateMaxMessageBytes(t *testing.T) {
	awsEnvironment := "cinp"
	supplierCode := "dev-2"

	topic := &client.Topic{
		Name: fmt.Sprintf("test_acc_%s", acctest.RandString(10)),
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckKafkaManagerTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckKafkaManagerTopicConfig_update_max_message_bytes(
					awsEnvironment,
					supplierCode,
					topic.Name,
					10240,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKafkaManagerTopicExists("kafkamanager_topic.foobar", topic),
					resource.TestCheckResourceAttr("kafkamanager_topic.foobar", "name", topic.Name),
					resource.TestCheckResourceAttr("kafkamanager_topic.foobar", "max_message_bytes", "10240"),
				),
			},
			{
				Config: testAccCheckKafkaManagerTopicConfig_update_max_message_bytes(
					awsEnvironment,
					supplierCode,
					topic.Name,
					20480,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKafkaManagerTopicExists("kafkamanager_topic.foobar", topic),
					resource.TestCheckResourceAttr("kafkamanager_topic.foobar", "name", topic.Name),
					resource.TestCheckResourceAttr("kafkamanager_topic.foobar", "max_message_bytes", "20480"),
				),
			},
		},
	})
}

func testAccCheckKafkaManagerTopicConfig_update_max_message_bytes(awsEnvironment string, supplierCode string, topicName string, maxMessageBytes int) string {
	return fmt.Sprintf(`
data "kafkamanager_cluster" "bazqux" {
	name = "data-platform-%s-%s-cluster"
}

resource "kafkamanager_topic" "foobar" {
	name               = "%s"
	cluster_id         = data.kafkamanager_cluster.bazqux.id
	partitions         = 6
	max_message_bytes  = %d
}`, awsEnvironment, supplierCode, topicName, maxMessageBytes)
}

func TestAccTopicModification_UpdateRetentionMs(t *testing.T) {
	awsEnvironment := "cinp"
	supplierCode := "dev-2"

	topic := &client.Topic{
		Name: fmt.Sprintf("test_acc_%s", acctest.RandString(10)),
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckKafkaManagerTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckKafkaManagerTopicConfig_update_retention_ms(
					awsEnvironment,
					supplierCode,
					topic.Name,
					3600,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKafkaManagerTopicExists("kafkamanager_topic.foobar", topic),
					resource.TestCheckResourceAttr("kafkamanager_topic.foobar", "name", topic.Name),
					resource.TestCheckResourceAttr("kafkamanager_topic.foobar", "retention_ms", "3600"),
				),
			},
			{
				Config: testAccCheckKafkaManagerTopicConfig_update_retention_ms(
					awsEnvironment,
					supplierCode,
					topic.Name,
					7200,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKafkaManagerTopicExists("kafkamanager_topic.foobar", topic),
					resource.TestCheckResourceAttr("kafkamanager_topic.foobar", "name", topic.Name),
					resource.TestCheckResourceAttr("kafkamanager_topic.foobar", "retention_ms", "7200"),
				),
			},
		},
	})
}

func testAccCheckKafkaManagerTopicConfig_update_retention_ms(awsEnvironment string, supplierCode string, topicName string, retentionMs int) string {
	return fmt.Sprintf(`
data "kafkamanager_cluster" "bazqux" {
	name = "data-platform-%s-%s-cluster"
}

resource "kafkamanager_topic" "foobar" {
	name               = "%s"
	cluster_id         = data.kafkamanager_cluster.bazqux.id
	partitions         = 6
	retention_ms	   = %d
}`, awsEnvironment, supplierCode, topicName, retentionMs)
}

func TestAccTopicModification_UpdateRetentionBytes(t *testing.T) {
	awsEnvironment := "cinp"
	supplierCode := "dev-2"

	topic := &client.Topic{
		Name: fmt.Sprintf("test_acc_%s", acctest.RandString(10)),
	}

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckKafkaManagerTopicDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckKafkaManagerTopicConfig_update_retention_bytes(
					awsEnvironment,
					supplierCode,
					topic.Name,
					1048576,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKafkaManagerTopicExists("kafkamanager_topic.foobar", topic),
					resource.TestCheckResourceAttr("kafkamanager_topic.foobar", "name", topic.Name),
					resource.TestCheckResourceAttr("kafkamanager_topic.foobar", "retention_bytes", "1048576"),
				),
			},
			{
				Config: testAccCheckKafkaManagerTopicConfig_update_retention_bytes(
					awsEnvironment,
					supplierCode,
					topic.Name,
					2097152,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckKafkaManagerTopicExists("kafkamanager_topic.foobar", topic),
					resource.TestCheckResourceAttr("kafkamanager_topic.foobar", "name", topic.Name),
					resource.TestCheckResourceAttr("kafkamanager_topic.foobar", "retention_bytes", "2097152"),
				),
			},
		},
	})
}

func testAccCheckKafkaManagerTopicConfig_update_retention_bytes(awsEnvironment string, supplierCode string, topicName string, retentionBytes int) string {
	return fmt.Sprintf(`
data "kafkamanager_cluster" "bazqux" {
	name = "data-platform-%s-%s-cluster"
}

resource "kafkamanager_topic" "foobar" {
	name               = "%s"
	cluster_id         = data.kafkamanager_cluster.bazqux.id
	partitions         = 6
	retention_bytes	   = %d
}`, awsEnvironment, supplierCode, topicName, retentionBytes)
}
