package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func setResourceDataFromMap(d *schema.ResourceData, m map[string]interface{}) error {
	for key, value := range m {
		if err := d.Set(key, value); err != nil {
			return fmt.Errorf("unable to set `%s` attribute: %s", key, err)
		}
	}
	return nil
}
