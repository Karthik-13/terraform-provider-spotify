package spotify

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zmb3/spotify"
)

type Range struct {
	Start int
	End   int
}

func resourcePlaylist() *schema.Resource {
	return &schema.Resource{
		Create: resourcePlaylistCreate,
		Read:   resourcePlaylistRead,
		Update: resourcePlaylistUpdate,
		Delete: resourcePlaylistDelete,

		Description: "Resource to manage a spotify playlist.",

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the playlist",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Description of playlist",
			},
			"tracks": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "Set of Track IDs for the playlist to contain",
			},
			"public": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Playlist should be accessible to public or not",
			},
			"snapshot_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourcePlaylistCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*spotify.Client)

	user, err := client.CurrentUser()
	if err != nil {
		return err
	}

	userID := string(user.ID)
	name := d.Get("name").(string)
	description := d.Get("description").(string)
	public := d.Get("public").(bool)

	playlist, err := client.CreatePlaylistForUser(userID, name, description, public)
	if err != nil {
		return fmt.Errorf("CreatePlaylist: %w", err)
	}

	d.SetId(string(playlist.ID))

	trackIDs := spotifyIdsInterface(d.Get("tracks").(*schema.Set).List())

	snapshotID := playlist.SnapshotID
	for _, rng := range batches(len(trackIDs), 100) {
		var err error
		snapshotID, err = client.AddTracksToPlaylist(playlist.ID, trackIDs[rng.Start:rng.End]...)
		if err != nil {
			return fmt.Errorf("AddTracksToPlaylist: %w", err)
		}
	}

	d.Set("snapshot_id", snapshotID)

	return nil
}

func resourcePlaylistRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*spotify.Client)

	playlistID := spotify.ID(d.Id())
	playlist, err := client.GetPlaylist(playlistID)

	if err != nil {
		return fmt.Errorf("GetPlaylist: %w", err)
	}

	d.Set("name", playlist.Name)
	d.Set("description", playlist.Description)
	d.Set("public", playlist.IsPublic)
	d.Set("snapshot_id", playlist.SnapshotID)

	trackIDs := schema.NewSet(schema.HashString, nil)

	tracks, err := client.GetPlaylistTracks(playlistID)
	if err != nil {
		return fmt.Errorf("GetPlaylistTracks: %w", err)
	}
	for err == nil {
		for _, track := range tracks.Tracks {
			trackIDs.Add(string(track.Track.ID))
		}
		err = client.NextPage(tracks)
	}

	d.Set("tracks", trackIDs)

	return nil
}

func resourcePlaylistUpdate(d *schema.ResourceData, m interface{}) error {
	client := m.(*spotify.Client)

	d.Partial(true)

	id := spotify.ID(d.Id())
	if d.HasChanges("name", "description", "public") {
		err := client.ChangePlaylistNameAccessAndDescription(
			id,
			d.Get("name").(string),
			d.Get("description").(string),
			d.Get("public").(bool),
		)

		if err != nil {
			return fmt.Errorf("ChangePlaylist: %w", err)
		}
	}

	if d.HasChange("tracks") {
		old, new := d.GetChange("tracks")
		oldSet := old.(*schema.Set)
		newSet := new.(*schema.Set)
		add := spotifyIdsInterface(newSet.Difference(oldSet).List())
		sub := spotifyIdsInterface(oldSet.Difference(newSet).List())

		var err error
		var snapshotID string
		for _, rng := range batches(len(add), 100) {
			snapshotID, err = client.AddTracksToPlaylist(id, add[rng.Start:rng.End]...)
			if err != nil {
				return fmt.Errorf("AddTracksToPlaylist: %w", err)
			}
		}
		for _, rng := range batches(len(sub), 100) {
			snapshotID, err = client.RemoveTracksFromPlaylist(id, sub[rng.Start:rng.End]...)
			if err != nil {
				return fmt.Errorf("AddTracksToPlaylist: %w", err)
			}
		}

		d.Set("snapshot_id", snapshotID)

	}

	return nil
}

func resourcePlaylistDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func spotifyIdsInterface(s []interface{}) []spotify.ID {
	output := make([]spotify.ID, len(s))
	for i, v := range s {
		output[i] = spotify.ID(v.(string))
	}
	return output
}

func batches(length, batch int) []Range {
	if length <= 0 {
		return nil
	}

	var ranges []Range
	i := 0
	for ; i < (length - batch); i += batch {
		ranges = append(ranges, Range{i, i + batch})
	}
	return append(ranges, Range{i, length})
}