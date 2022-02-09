package provider

import (
	"context"
	"strconv"
	"time"

	"coxautoinc.com/data-platform/kafka-manager/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceSchemaRegistries() *schema.Resource {
	recordSchema := schemaRegistrySchema()

	for _, f := range recordSchema {
		f.Computed = true
	}

	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"schema_registries": &schema.Schema{
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: recordSchema,
				},
			},
		},
		ReadContext: dataSourceSchemaRegistriesRead,
	}
}

func dataSourceSchemaRegistriesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(client.Client)

	rawSchemaRegistries, err := client.GetSchemaRegistries()
	if err != nil {
		return diag.FromErr(err)
	}

	schemaRegistries, err := marshalSchemaRegistries(&rawSchemaRegistries)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("schema_registries", schemaRegistries); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}

func marshalSchemaRegistries(schemaRegistries *[]client.SchemaRegistry) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)

	for _, schemaRegistry := range *schemaRegistries {
		sr, err := marshalSchemaRegistry(&schemaRegistry)
		if err != nil {
			return nil, err
		}
		result = append(result, sr)
	}

	return result, nil
}
