package spotify

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zmb3/spotify"
)

func resourcePlayerVolume() *schema.Resource {
	return &schema.Resource{
		Create: resourcePlayerVolumeCreate,
		Update: resourcePlayerVolumeCreate,
		Delete: resourcePlayerVolumeDelete,
		Read:   resourcePlayerVolumeRead,

		Description: "Resource to manage spotify player's volume on user's active device.",

		Schema: map[string]*schema.Schema{
			"percent": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The percent of player volume",
			},
		},
	}
}

func resourcePlayerVolumeCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*spotify.Client)

	err := client.Volume(d.Get("percent").(int))
	if err != nil {
		return err
	}

	return nil
}

func resourcePlayerVolumeDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourcePlayerVolumeRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

