package provider

import (
	"context"
	"strconv"
	"time"

	"coxautoinc.com/data-platform/kafka-manager/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceClusters() *schema.Resource {
	recordSchema := clusterSchema()

	for _, f := range recordSchema {
		f.Computed = true
	}

	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"clusters": &schema.Schema{
				Computed: true,
				Type:     schema.TypeList,
				Elem: &schema.Resource{
					Schema: recordSchema,
				},
			},
		},
		ReadContext: dataSourceClustersRead,
	}
}

func dataSourceClustersRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(client.Client)

	rawClusters, err := c.GetClusters()
	if err != nil {
		return diag.FromErr(err)
	}

	clusters, err := marshalClusters(&rawClusters)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("clusters", clusters); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return nil
}

func marshalClusters(clusters *[]client.Cluster) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)

	for _, cluster := range *clusters {
		c, err := marshalCluster(&cluster)
		if err != nil {
			return nil, err
		}
		result = append(result, c)
	}

	return result, nil
}
