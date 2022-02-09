package provider

import (
	"context"
	"strconv"
	"time"

	"coxautoinc.com/data-platform/kafka-manager/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceEnvironments() *schema.Resource {
	recordSchema := environmentSchema()

	for _, f := range recordSchema {
		f.Computed = true
	}

	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"environments": &schema.Schema{
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: recordSchema,
				},
			},
		},
		ReadContext: dataSourceEnvironmentsRead,
	}
}

func dataSourceEnvironmentsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Client)

	rawEnvironments, err := c.GetEnvironments()
	if err != nil {
		return diag.FromErr(err)
	}

	environments, err := marshalEnvironments(&rawEnvironments)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("environments", environments); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}

func marshalEnvironments(environments *[]client.Environment) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)

	for _, environment := range *environments {
		e, err := marshalEnvironment(&environment)
		if err != nil {
			return nil, err
		}
		result = append(result, e)
	}

	return result, nil
}
