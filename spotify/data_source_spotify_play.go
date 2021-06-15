package spotify

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zmb3/spotify"
)

func dataSourceSpotifyPlay() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSpotifyPlayRead,
		Schema: map[string]*schema.Schema{
		},
	}
}

func dataSourceSpotifyPlayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*spotify.Client)

	err := client.Play()
	if err != nil {
		return err
	}

	return nil
}