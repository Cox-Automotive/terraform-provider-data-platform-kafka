package provider

import (
	"context"
	"fmt"
	"strconv"

	"coxautoinc.com/data-platform/kafka-manager/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func topicSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": &schema.Schema{
			Type: schema.TypeString,
		},
		"cluster_id": &schema.Schema{
			Type: schema.TypeString,
		},
		"name": &schema.Schema{
			Type: schema.TypeString,
		},
		"partitions": &schema.Schema{
			Type: schema.TypeInt,
		},
		"replication_factor": &schema.Schema{
			Type: schema.TypeInt,
		},
		"max_message_bytes": &schema.Schema{
			Type: schema.TypeInt,
		},
		"cleanup_policy": &schema.Schema{
			Type: schema.TypeString,
		},
		"retention_ms": &schema.Schema{
			Type: schema.TypeInt,
		},
		"retention_bytes": &schema.Schema{
			Type: schema.TypeInt,
		},
	}
}

func resourceTopic() *schema.Resource {
	recordSchema := topicSchema()
	recordSchema["id"].Computed = true
	recordSchema["name"].Required = true
	recordSchema["cluster_id"].Required = true
	recordSchema["partitions"].Optional = true
	recordSchema["partitions"].Computed = true
	recordSchema["max_message_bytes"].Optional = true
	recordSchema["max_message_bytes"].Computed = true
	recordSchema["cleanup_policy"].Optional = true
	recordSchema["cleanup_policy"].Computed = true
	recordSchema["retention_ms"].Optional = true
	recordSchema["retention_ms"].Computed = true
	recordSchema["retention_bytes"].Optional = true
	recordSchema["retention_bytes"].Computed = true
	recordSchema["replication_factor"].Optional = true
	recordSchema["replication_factor"].Computed = true

	return &schema.Resource{
		Schema:        recordSchema,
		CreateContext: resourceTopicCreate,
		ReadContext:   resourceTopicRead,
		UpdateContext: resourceTopicUpdate,
		DeleteContext: resourceTopicDelete,
	}
}

func marshalTopic(t *client.Topic) (map[string]interface{}, error) {
	result := map[string]interface{}{
		"id":                 strconv.Itoa(t.ID),
		"name":               t.Name,
		"cluster_id":         strconv.Itoa(t.Cluster.ID),
		"partitions":         t.Partitions,
		"max_message_bytes":  t.Config.MaxMessageBytes,
		"cleanup_policy":     t.Config.CleanupPolicy,
		"retention_ms":       t.Config.RetentionMs,
		"retention_bytes":    t.Config.RetentionBytes,
		"replication_factor": t.Config.ReplicationFactor,
	}

	return result, nil
}

func unmarshalNewTopic(d *schema.ResourceData) (*client.NewTopic, error) {
	clusterID, err := strconv.Atoi(d.Get("cluster_id").(string))
	if err != nil {
		return nil, fmt.Errorf("invalid cluster ID: %s", err)
	}

	topic := &client.NewTopic{
		Name:      d.Get("name").(string),
		ClusterID: clusterID,
		Config:    &client.TopicConfig{},
	}

	if v, ok := d.GetOk("partitions"); ok {
		topic.Partitions = v.(int)
	}
	if v, ok := d.GetOk("replication_factor"); ok {
		topic.Config.ReplicationFactor = v.(int)
	}
	if v, ok := d.GetOk("max_message_bytes"); ok {
		topic.Config.MaxMessageBytes = v.(int)
	}
	if v, ok := d.GetOk("cleanup_policy"); ok {
		topic.Config.CleanupPolicy = v.(string)
	}
	if v, ok := d.GetOk("retention_ms"); ok {
		topic.Config.RetentionMs = v.(int)
	}
	if v, ok := d.GetOk("retention_bytes"); ok {
		topic.Config.RetentionBytes = v.(int)
	}

	return topic, nil
}

func resourceTopicCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(client.Client)
	newTopic, err := unmarshalNewTopic(d)
	if err != nil {
		return diag.FromErr(err)
	}

	topic, err := c.CreateTopic(newTopic)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(topic.ID))
	return resourceTopicRead(ctx, d, meta)
}

func resourceTopicRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(client.Client)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("invalid topic ID: %s", err)
	}

	rawTopic, err := c.GetTopic(id)
	if err != nil {
		return diag.Errorf("error reading topic: %s", err)
	}

	topic, err := marshalTopic(rawTopic)
	if err != nil {
		return diag.Errorf("error reading topic: %s", err)
	}
	d.SetId(topic["id"].(string))
	setResourceDataFromMap(d, topic)

	return nil
}

func resourceTopicUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(client.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("invalid topic ID: %s", err)
	}

	topic := &client.Topic{
		ID:     id,
		Config: &client.TopicConfig{},
	}

	if d.HasChange("max_message_bytes") {
		topic.Config.MaxMessageBytes = d.Get("max_message_bytes").(int)
	}
	if d.HasChange("cleanup_policy") {
		topic.Config.CleanupPolicy = d.Get("cleanup_policy").(string)
	}
	if d.HasChange("retention_ms") {
		topic.Config.RetentionMs = d.Get("retention_ms").(int)
	}
	if d.HasChange("retention_bytes") {
		topic.Config.RetentionBytes = d.Get("retention_bytes").(int)
	}

	err = c.UpdateTopic(topic)
	if err != nil {
		return diag.Errorf("failed to update topic: %s", err)
	}

	return resourceTopicRead(ctx, d, meta)
}

func resourceTopicDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(client.Client)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.Errorf("invalid topic ID: %s", err)
	}

	err = c.DeleteTopic(id)
	if err != nil {
		return diag.Errorf("failed to delete topic: %s", err)
	}

	return nil
}
