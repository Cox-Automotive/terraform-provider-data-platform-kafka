package provider

import (
	"context"
	"strconv"

	"coxautoinc.com/data-platform/kafka-manager/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func schemaRegistrySchema() map[string]*schema.Schema {
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
		"endpoint": &schema.Schema{
			Type: schema.TypeString,
		},
		"cloud": &schema.Schema{
			Type: schema.TypeString,
		},
		"region": &schema.Schema{
			Type: schema.TypeString,
		},
		"secrets_manager_secret_arn": &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

func dataSourceSchemaRegistry() *schema.Resource {
	recordSchema := schemaRegistrySchema()

	for _, f := range recordSchema {
		f.Computed = true
	}

	recordSchema["id"].ExactlyOneOf = []string{"id", "confluent_id"}
	recordSchema["id"].Optional = true
	recordSchema["confluent_id"].ExactlyOneOf = []string{"id", "confluent_id"}
	recordSchema["confluent_id"].Optional = true

	return &schema.Resource{
		Schema:      recordSchema,
		ReadContext: dataSourceSchemaRegistryRead,
	}
}

func dataSourceSchemaRegistryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Client)

	var rawSchemaRegistry *client.SchemaRegistry
	var err error

	if id, ok := d.GetOk("id"); ok {
		id, err = strconv.Atoi(id.(string))
		if err != nil {
			return diag.Errorf("invalid schema registry ID: %s", err)
		}
		rawSchemaRegistry, err = c.GetSchemaRegistry(id.(int))
		if err != nil {
			return diag.FromErr(err)
		}
	} else if confluentID, ok := d.GetOk("confluent_id"); ok {
		rawSchemaRegistry, err = c.GetSchemaRegistryByConfluentID(confluentID.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		return diag.Errorf("provide either schema registry id or confluent_id")
	}

	schemaRegistry, err := marshalSchemaRegistry(rawSchemaRegistry)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := setResourceDataFromMap(d, schemaRegistry); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(schemaRegistry["id"].(string))

	return nil
}

func marshalSchemaRegistry(sr *client.SchemaRegistry) (map[string]interface{}, error) {
	result := map[string]interface{}{
		"id":                         strconv.Itoa(sr.ID),
		"environment_id":             sr.Environment.ID,
		"confluent_id":               sr.ConfluentID,
		"endpoint":                   sr.Endpoint,
		"cloud":                      sr.Cloud,
		"region":                     sr.Region,
		"secrets_manager_secret_arn": sr.SecretsManagerSecretArn,
	}

	return result, nil
}
