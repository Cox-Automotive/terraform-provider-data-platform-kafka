package provider

import (
	"context"
	"strconv"

	"coxautoinc.com/data-platform/kafka-manager/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func environmentSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": &schema.Schema{
			Type: schema.TypeString,
		},
		"name": &schema.Schema{
			Type: schema.TypeString,
		},
		"confluent_id": &schema.Schema{
			Type: schema.TypeString,
		},
		"supplier": &schema.Schema{
			Type: schema.TypeString,
		},
	}
}

func dataSourceEnvironment() *schema.Resource {
	recordSchema := environmentSchema()

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
		ReadContext: dataSourceEnvironmentRead,
	}
}

func dataSourceEnvironmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Client)

	var rawEnvironment *client.Environment
	var err error

	if id, ok := d.GetOk("id"); ok {
		id, err = strconv.Atoi(id.(string))
		if err != nil {
			return diag.Errorf("invalid environment ID: %s", err)
		}
		rawEnvironment, err = c.GetEnvironment(id.(int))
		if err != nil {
			return diag.FromErr(err)
		}
	} else if confluentID, ok := d.GetOk("confluent_id"); ok {
		rawEnvironment, err = c.GetEnvironmentByConfluentID(confluentID.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	} else if name, ok := d.GetOk("name"); ok {
		rawEnvironment, err = c.GetEnvironmentByName(name.(string))
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		return diag.Errorf("provide either environment id, confluent_id, or name")
	}

	environment, err := marshalEnvironment(rawEnvironment)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := setResourceDataFromMap(d, environment); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(environment["id"].(string))

	return nil
}

func marshalEnvironment(e *client.Environment) (map[string]interface{}, error) {
	result := map[string]interface{}{
		"id":           strconv.Itoa(e.ID),
		"name":         e.Name,
		"confluent_id": e.ConfluentID,
		"supplier":     e.Supplier,
	}

	return result, nil
}
