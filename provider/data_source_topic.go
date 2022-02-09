package provider

import (
	"context"
	"strconv"

	"coxautoinc.com/data-platform/kafka-manager/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceTopic() *schema.Resource {
	recordSchema := topicSchema()

	for _, f := range recordSchema {
		f.Computed = true
	}

	recordSchema["id"].ExactlyOneOf = []string{"id", "name"}
	recordSchema["id"].Optional = true
	recordSchema["name"].ExactlyOneOf = []string{"id", "name"}
	recordSchema["name"].Optional = true
	recordSchema["cluster_id"].RequiredWith = []string{"name"}
	recordSchema["cluster_id"].Optional = true

	return &schema.Resource{
		Schema:      recordSchema,
		ReadContext: dataSourceTopicRead,
	}
}

func dataSourceTopicRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(client.Client)

	var rawTopic *client.Topic
	var err error

	if id, ok := d.GetOk("id"); ok {
		id, err = strconv.Atoi(id.(string))
		if err != nil {
			return diag.Errorf("invalid topic ID: %s", err)
		}
		rawTopic, err = c.GetTopic(id.(int))
		if err != nil {
			return diag.FromErr(err)
		}
	} else if name, ok := d.GetOk("name"); ok {
		if clusterID, ok := d.GetOk("cluster_id"); ok {
			clusterID, err = strconv.Atoi(clusterID.(string))
			if err != nil {
				return diag.Errorf("invalid cluster ID: %s", err)
			}
			rawTopic, err = c.GetTopicByNameAndClusterID(name.(string), clusterID.(int))
			if err != nil {
				return diag.FromErr(err)
			}
		} else {
			return diag.Errorf("provide both topic name and cluster id")
		}
	} else {
		return diag.Errorf("provide topic id, or topic name and cluster id")
	}

	topic, err := marshalTopic(rawTopic)
	if err != nil {
		return diag.FromErr(err)
	}
	if err := setResourceDataFromMap(d, topic); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(topic["id"].(string))

	return nil
}
