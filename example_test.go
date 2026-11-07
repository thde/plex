package plex_test

import (
	"context"
	"fmt"
	"log"

	"thde.io/plex"
)

func ExampleNew() {
	client, err := plex.New(
		"http://192.168.1.10:32400",
		plex.WithToken("your-plex-token"),
		plex.WithClientIdentifier("a-unique-stable-client-id"),
		plex.WithProduct("my-app", "1.0"),
	)
	if err != nil {
		log.Fatal(err)
	}

	resp, err := client.GetIdentityWithResponse(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	if resp.JSON200 == nil {
		log.Fatalf("unexpected response: %s", resp.Status())
	}

	fmt.Println(resp.JSON200.MediaContainer.Version)
}
