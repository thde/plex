// Command plex-openapi extracts the Plex Media Server OpenAPI specification
// from the Redoc documentation page at https://developer.plex.tv/pms/ and
// prints it as JSON.
//
// Plex does not publish the specification as a standalone file. The
// documentation page is a pre-rendered Redoc page that embeds the full
// specification in its hydration state (__redoc_state.spec.data).
//
// Usage:
//
//	plex-openapi [-url URL] [-o FILE]
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"
)

const defaultURL = "https://developer.plex.tv/pms/"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "plex-openapi:", err)
		os.Exit(1)
	}
}

func run() error {
	var (
		url     = flag.String("url", defaultURL, "URL of the Plex Media Server API documentation page")
		output  = flag.String("o", "", "write the specification to `file` instead of stdout")
		timeout = flag.Duration("timeout", time.Minute, "HTTP request timeout")
	)
	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [flags]\n\nExtracts the OpenAPI specification embedded in the Plex API documentation page.\n\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() > 0 {
		flag.Usage()
		return errors.New("unexpected arguments")
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	ctx, cancel := context.WithTimeout(ctx, *timeout)
	defer cancel()

	spec, err := fetchSpec(ctx, *url)
	if err != nil {
		return err
	}

	if *output == "" {
		_, err = os.Stdout.Write(spec)
		return err
	}
	return os.WriteFile(*output, spec, 0o644)
}
