package provider

import (
	"fmt"
	"os"
	"testing"

	"coxautoinc.com/data-platform/kafka-manager/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestMain(m *testing.M) {
	resource.TestMain(m)
}

func sharedPrivateClientForRegion(region string) (*client.PrivateClient, error) {
	// NOTE: region is not used at this moment and is only needed to conform to Terraform testing API.

	url := os.Getenv("KAFKAMANAGER_URL")
	if url == "" {
		url = "https://edp-microservices.awsedscinp.com/kafka-manager"
	}

	key := os.Getenv("KAFKAMANAGER_KEY")

	oktaGroups := os.Getenv("KAFKAMANAGER_OKTA_GROUPS")
	if oktaGroups == "" {
		oktaGroups = "EDP_DEV_OPS_DELETE"
	}

	supplier := os.Getenv("KAFKAMANAGER_SUPPLIER")
	if supplier == "" {
		supplier = "dev-1"
	}

	userID := os.Getenv("KAFKAMANAGER_USER_ID")
	if userID == "" {
		userID = "EDP_INTERNAL_APPID"
	}

	if key != "" {
		return client.NewPrivateClient(url, key, oktaGroups, supplier, userID), nil
	} else {
		return nil, fmt.Errorf("provide KAFKAMANAGER_KEY")
	}
}
