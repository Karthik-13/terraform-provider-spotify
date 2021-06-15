package spotify

import (
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zmb3/spotify"
)

func dataSourceSpotifyPlayerDevices() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSpotifyPlayerDevicesRead,
		Schema: map[string]*schema.Schema{
			"devices": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "",
						},
						"is_active": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "",
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
							Description: "",
						},
						"volume_percent": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "",
						},
					},
				},
			},
		},
	}
}

func dataSourceSpotifyPlayerDevicesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*spotify.Client)

	result, err := client.PlayerDevices()
	if err != nil {
		return err
	}

	var devices []interface{}

    for _, device := range result {
		devices = append(devices, map[string]interface{}{
			"name": device.Name,
		    "is_active": device.Active,
		    "type": device.Type,
		    "volume_percent": device.Volume,
		    "id": device.ID,
		})
	}

	d.Set("devices", devices)
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))
	return nil
}