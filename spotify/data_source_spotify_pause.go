package spotify

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zmb3/spotify"
)

func dataSourceSpotifyPause() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSpotifyPauseRead,
		Schema: map[string]*schema.Schema{
		},
	}
}

func dataSourceSpotifyPauseRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*spotify.Client)

	err := client.Pause()
	if err != nil {
		return err
	}

	return nil
}