package provider

import (
	"context"
	"strconv"
	"time"

	"coxautoinc.com/data-platform/kafka-manager/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTopics() *schema.Resource {
	recordSchema := topicSchema()

	for _, f := range recordSchema {
		f.Computed = true
	}

	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"topics": &schema.Schema{
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: recordSchema,
				},
			},
		},
		ReadContext: dataSourceTopicsRead,
	}
}

func dataSourceTopicsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(client.Client)

	rawTopics, err := client.GetTopics()
	if err != nil {
		return diag.FromErr(err)
	}

	topics, err := marshalTopics(&rawTopics)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("topics", topics); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}

func marshalTopics(topics *[]client.Topic) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)

	for _, topic := range *topics {
		t, err := marshalTopic(&topic)
		if err != nil {
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}
