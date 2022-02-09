package provider

import (
	"context"
	"strconv"

	"coxautoinc.com/data-platform/kafka-manager/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func clusterSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": &schema.Schema{
			Type: schema.TypeString,
		},
		"environment_id": &schema.Schema{
			Type: schema.TypeInt,
		},
		"confluent_id": &schema.Schema{
			Type: schema.TypeString,
		},
		"name": &schema.Schema{
			Type: schema.TypeString,
		},
		"cloud": &schema.Schema{
			Type: schema.TypeString,
		},
		"region": &schema.Schema{
			Type: schema.TypeString,
		},
		"availability": &schema.Schema{
			Type: schema.TypeString,
		},
		"cluster_type": &schema.Schema{
			Type: schema.TypeString,
		},
		"networking_type": &schema.Schema{
			Type: schema.TypeString,
		},
		"secrets_manager_secret_arn": &schema.Schema{
			Type: schema.TypeString,
		},
		"secrets_manager_cloud_secret_arn": &schema.Schema{
			Type: schema.TypeString,
		},
		"rest_proxy_url": &schema.Schema{
			Type: schema.TypeString,
		},
		"bootstrap_servers": &schema.Schema{
			Type: schema.TypeString,
		},
		"default_s3_bucket_name": &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

func dataSourceCluster() *schema.Resource {
	recordSchema := clusterSchema()

	for _, f := range recordSchema {
		f.Computed = true
	}

	recordSchema["id"].ExactlyOneOf = []string{"id", "confluent_id", "name"}
	recordSchema["id"].Optional = true
	recordSchema["confluent_id"].ExactlyOneOf = []string{"id", "confluent_id", "name"}
	recordSchema["confluent_id"].Optional = true
	recordSchema["name"].ExactlyOneOf = []string{"id", "confluent_id", "name"}
	recordSchema["name"].Optional = true

	return &schema.Resource{
		Schema:      recordSchema,
		ReadContext: dataSourceClusterRead,
	}
}

func dataSourceClusterRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(client.Client)

	var rawCluster *client.Cluster
	var err error

	if id, ok := d.GetOk("id"); ok {
		id, err = strconv.Atoi(id.(string))
		if err != nil {
			return diag.Errorf("invalid cluster ID: %s", err)
		}
		rawCluster, err = c.GetCluster(id.(int))
		if err != nil {
			return diag.FromErr(err)
		}
	} else if confluentID, ok := d.GetOk("confluent_id"); ok {
		rawCluster, err = c.GetClusterByConfluentID(confluentID.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	} else if name, ok := d.GetOk("name"); ok {
		rawCluster, err = c.GetClusterByName(name.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		return diag.Errorf("provide either cluster id, confluent_id, or name")
	}

	cluster, err := marshalCluster(rawCluster)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := setResourceDataFromMap(d, cluster); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(cluster["id"].(string))

	return nil
}

func marshalCluster(c *client.Cluster) (map[string]interface{}, error) {
	result := map[string]interface{}{
		"id":                               strconv.Itoa(c.ID),
		"name":                             c.Name,
		"confluent_id":                     c.ConfluentID,
		"environment_id":                   c.Environment.ID,
		"cloud":                            c.Cloud,
		"region":                           c.Region,
		"availability":                     c.Availability,
		"cluster_type":                     c.ClusterType,
		"networking_type":                  c.NetworkingType,
		"secrets_manager_secret_arn":       c.SecretsManagerSecretArn,
		"secrets_manager_cloud_secret_arn": c.SecretsManagerCloudSecretArn,
		"rest_proxy_url":                   c.RestProxyUrl,
		"bootstrap_servers":                c.BootstrapServers,
		"default_s3_bucket_name":           c.DefaultS3BucketName,
	}

	return result, nil
}
