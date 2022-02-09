package provider

import (
	"context"

	"coxautoinc.com/data-platform/kafka-manager/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("KAFKAMANAGER_URL", nil),
			},
			"access_token": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"key", "okta_groups", "supplier", "user_id"},
				ExactlyOneOf:  []string{"access_token", "key"},
				DefaultFunc:   schema.EnvDefaultFunc("KAFKAMANAGER_ACCESS_TOKEN", nil),
			},
			"key": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				ConflictsWith: []string{"access_token"},
				ExactlyOneOf:  []string{"access_token", "key"},
				DefaultFunc:   schema.EnvDefaultFunc("KAFKAMANAGER_KEY", nil),
			},
			"okta_groups": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"access_token"},
				RequiredWith:  []string{"key"},
				DefaultFunc:   schema.EnvDefaultFunc("KAFKAMANAGER_OKTA_GROUPS", nil),
			},
			"supplier": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"access_token"},
				RequiredWith:  []string{"key"},
				DefaultFunc:   schema.EnvDefaultFunc("KAFKAMANAGER_SUPPLIER", nil),
			},
			"user_id": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"access_token"},
				RequiredWith:  []string{"key"},
				DefaultFunc:   schema.EnvDefaultFunc("KAFKAMANAGER_USER_ID", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"kafkamanager_topic": resourceTopic(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"kafkamanager_environment":       dataSourceEnvironment(),
			"kafkamanager_environments":      dataSourceEnvironments(),
			"kafkamanager_schema_registry":   dataSourceSchemaRegistry(),
			"kafkamanager_schema_registries": dataSourceSchemaRegistries(),
			"kafkamanager_cluster":           dataSourceCluster(),
			"kafkamanager_clusters":          dataSourceClusters(),
			"kafkamanager_topic":             dataSourceTopic(),
			"kafkamanager_topics":            dataSourceTopics(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	url := d.Get("url").(string)

	if accessToken, ok := d.GetOk("access_token"); ok {
		return client.NewPublicClient(url, accessToken.(string)), nil
	} else if key, ok := d.GetOk("key"); ok {
		oktaGroups := d.Get("okta_groups").(string)
		supplier := d.Get("supplier").(string)
		userID := d.Get("user_id").(string)
		return client.NewPrivateClient(url, key.(string), oktaGroups, supplier, userID), nil
	} else {
		return nil, diag.Errorf("provide either access token or Data-Platform Key, Okta groups, supplier code, and User ID")
	}
}
