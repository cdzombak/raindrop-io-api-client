# Unofficial Raindrop.io API client
## Work in progress

[![Actions Status](https://github.com/kattaris/raindrop-io-api-client/workflows/CI/badge.svg)](https://github.com/kattaris/raindrop-io-api-client/actions)
[![Coverage Status](https://codecov.io/github/kattaris/raindrop-io-api-client/coverage.svg?branch=master)](https://codecov.io/gh/kattaris/raindrop-io-api-client)
[![Releases](https://img.shields.io/github/v/release/kattaris/raindrop-io-api-client.svg?include_prereleases&style=flat-square)](https://github.com/kattaris/raindrop-io-api-client/releases)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

### Example usage:

```
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
	// Create logger (optional)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	
	client, err := raindrop.NewClientWithLogger("5478394jfkdlsf843u430",
		"e46b6a8a-018d-43b1-8b28-543kjl32ghj",
		"http://localhost:8080/oauth", logger)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		http.HandleFunc("/oauth", client.GetAuthorizationCodeHandler)
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	}()

	// Step 1: The authorization request
	authUrl, err := client.GetAuthorizationURL()
	if err != nil {
		log.Fatal(err)
	}

	// Step 2: The redirection to your application site
	u, err := url.QueryUnescape(authUrl.String())
	if err != nil {
		log.Fatal(err)
	}
	if err := openBrowser(u); err != nil {
		log.Fatal(err)
	}

	// Step 3: The token exchange
	for client.ClientCode == "" {
		fmt.Println("Waiting for client to authorize")
		time.Sleep(3 * time.Second)
	}

	ctx := context.Background()
	accessTokenResp, err := client.GetAccessToken(client.ClientCode, ctx)
	if err != nil {
		log.Fatal(err)
	}
	accessToken := accessTokenResp.AccessToken

	// Step 4: Check API's methods
	result, err := client.CreateCollection(accessToken, true, "list",
		"Test", 1, false, 0, nil, ctx)
	if err != nil {
		log.Printf("Error creating collection: %v", err)
	} else {
		fmt.Printf("Create collection result: %v\n", result)
	}

	rootCollections, err := client.GetRootCollections(accessToken, ctx)
	if err != nil {
		log.Printf("Error getting root collections: %v", err)
	} else {
		fmt.Printf("Root Collections: %v\n", rootCollections)
	}

	childCollections, err := client.GetChildCollections(accessToken, ctx)
	if err != nil {
		log.Printf("Error getting child collections: %v", err)
	} else {
		fmt.Printf("Child Collections: %v\n", childCollections)
	}
}

func openBrowser(url string) error {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	return err
}
```
