package spotify

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zmb3/spotify"
)

func dataSourceSpotifyUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSpotifyUserRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
				Description: "ID of the User",
			},
			"display_name": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "Name of the User",
			},
			"endpoint": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "Endpoint of the User",
			},
			"external_urls": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "External URL of the User",
			},
			"uri": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "URI of the User",
			},
		},
	}
}

func dataSourceSpotifyUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*spotify.Client)

	user, err := client.GetUsersPublicProfile(spotify.ID(d.Get("id").(string)))
	if err != nil {
		return err
	}

	d.Set("external_urls", user.ExternalURLs)
	d.Set("display_name", user.DisplayName)
	d.Set("endpoint", user.Endpoint)
	d.Set("uri", user.URI)
	d.SetId(string(user.ID))

	return nil
}