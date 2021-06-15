package spotify

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zmb3/spotify"
)

func dataSourceSpotifyAlbum() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSpotifyAlbumRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
				Description: "ID of the Spotify Album",
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "The Name of the Album",
			},
		},
	}
}

func dataSourceSpotifyAlbumRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*spotify.Client)

	album, err := client.GetAlbum(spotify.ID(d.Get("spotify_id").(string)))
	if err != nil {
		return err
	}

	d.Set("name", album.Name)
	d.SetId(string(album.ID))

	return nil
}