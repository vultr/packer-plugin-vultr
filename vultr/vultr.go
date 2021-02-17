package vultr

import (
	"context"
	"fmt"

	"github.com/hashicorp/packer-plugin-sdk/version"
	"github.com/vultr/govultr/v2"
	"golang.org/x/oauth2"
)

func newVultrClient(apiKey string) *govultr.Client {
	ctx := context.Background()

	config := &oauth2.Config{}
	ts := config.TokenSource(ctx, &oauth2.Token{AccessToken: apiKey})

	client := govultr.NewClient(oauth2.NewClient(ctx, ts))
	userAgent := fmt.Sprintf("Packer/%s/govultr-v2", version.SDKVersion)
	client.SetUserAgent(userAgent)
	return client
}
