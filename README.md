# Unofficial Raindrop.io API client

A Go client for the [Raindrop.io](https://raindrop.io) [REST API](https://developer.raindrop.io/). It handles the OAuth 2.0 authorization flow and wraps the collection, raindrop, and tag endpoints.

## Installation

```shell
go get github.com/cdzombak/raindrop-io-api-client
```

Requires the Go version declared in [`go.mod`](go.mod).

## Usage

Create a client with your Raindrop.io app credentials, complete the OAuth flow to obtain an access token, then call the API methods. Logging is optional; use `NewClient` to omit it or `NewClientWithLogger` to pass an `slog.Logger`.

```go
package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/cdzombak/raindrop-io-api-client/pkg/raindrop"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	client, err := raindrop.NewClientWithLogger(
		"your-client-id",
		"your-client-secret",
		"http://localhost:8080/oauth",
		logger,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Serve the OAuth redirect endpoint.
	go func() {
		http.HandleFunc("/oauth", client.GetAuthorizationCodeHandler)
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	}()

	// Send the user to Raindrop.io to authorize the app.
	authURL, err := client.GetAuthorizationURL()
	if err != nil {
		log.Fatal(err)
	}
	u, err := url.QueryUnescape(authURL.String())
	if err != nil {
		log.Fatal(err)
	}
	if err := openBrowser(u); err != nil {
		log.Fatal(err)
	}

	// Wait for the authorization code, then exchange it for an access token.
	for client.ClientCode == "" {
		fmt.Println("Waiting for authorization…")
		time.Sleep(3 * time.Second)
	}

	ctx := context.Background()
	tokenResp, err := client.GetAccessToken(client.ClientCode, ctx)
	if err != nil {
		log.Fatal(err)
	}
	accessToken := tokenResp.AccessToken

	// Call the API.
	collections, err := client.GetRootCollections(accessToken, ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Root collections: %v\n", collections)
}

func openBrowser(url string) error {
	switch runtime.GOOS {
	case "linux":
		return exec.Command("xdg-open", url).Start()
	case "windows":
		return exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default:
		return fmt.Errorf("unsupported platform")
	}
}
```

Once you have an access token, you can refresh it later with `RefreshAccessToken` rather than repeating the full authorization flow.

## Supported endpoints

- **Authorization:** `GetAuthorizationURL`, `GetAuthorizationCodeHandler`, `GetAuthorizationCode`, `GetAccessToken`, `RefreshAccessToken`
- **Collections:** `GetRootCollections`, `GetChildCollections`, `GetCollection`, `CreateCollection`
- **Raindrops:** `CreateSimpleRaindrop`, `GetRaindrops`, `GetTaggedRaindrops`
- **Tags:** `GetTags`, `DeleteTags`

## License

Apache 2.0. See [LICENSE](LICENSE).

This is a fork of [antonnagorniy/raindrop-io-api-client](https://github.com/antonnagorniy/raindrop-io-api-client).
