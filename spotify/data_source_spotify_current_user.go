package spotify

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zmb3/spotify"
)

func dataSourceSpotifyCurrentUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSpotifyCurrentUserRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
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
			"email": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "External URL of the User",
			},
			"subscription": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "Subscription of the User",
			},
			"country": {
				Type:     schema.TypeString,
				Computed: true,
				Description: "Country of the User",
			},
		},
	}
}

func dataSourceSpotifyCurrentUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*spotify.Client)

	user, err := client.CurrentUser()
	if err != nil {
		return err
	}

	d.Set("email", user.Email)
	d.Set("display_name", user.DisplayName)
	d.Set("endpoint", user.Endpoint)
	d.Set("subscription", user.Product)
	d.Set("country", user.Country)
	d.SetId(string(user.ID))

	return nil
}