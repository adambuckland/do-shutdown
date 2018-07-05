package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/digitalocean/godo"

	"golang.org/x/oauth2"
)

const accessTokenVariable = "DIGITALOCEAN_ACCESS_TOKEN"

var (
	dryrun = flag.Bool("dryrun", false, "Will do a dry run, but not actually delete anything")
	token  = flag.String("token", os.Getenv(accessTokenVariable), "DigitalOcean Personal Access Token")
)

func main() {
	flag.Parse()
	if *token == "" {
		fmt.Printf("No token specified. Please specifiy either with -token flag or %s env variable\n", accessTokenVariable)
		os.Exit(1)
	}

	oauthClient := oauth2.NewClient(context.Background(), &tokenSource{AccessToken: *token})
	client := godo.NewClient(oauthClient)

	droplets, err := fetchDroplets(client)
	if err != nil {
		fmt.Printf("❌ could not fetch droplets: %v\n", err)
		os.Exit(1)
	}

	success := true

	if len(droplets) == 0 {
		fmt.Printf("No droplets! All done 🎉\n")
	} else {
		fmt.Printf("-> 💣 Deleting %d droplet(s)\n", len(droplets))
		for _, droplet := range droplets {
			fmt.Printf("\t-> Droplet %s...", droplet.Name)
			if !*dryrun {
				_, err := client.Droplets.Delete(context.TODO(), droplet.ID)
				if err != nil {
					fmt.Printf("ERR: %v\n", err)
					success = false
					continue
				}
				fmt.Printf("\t DELETED\n")
			} else {
				fmt.Printf("\t (not) DELETED\n")
			}
		}
	}
	if !success {
		fmt.Printf("😞 Could not delete all droplets. Please see individual errors\n")
	}
}

func fetchDroplets(client *godo.Client) ([]godo.Droplet, error) {
	droplets := []godo.Droplet{}
	opt := &godo.ListOptions{}
	for {
		d, rsp, err := client.Droplets.List(context.TODO(), opt)
		if err != nil {
			return droplets, err
		}

		droplets = append(droplets, d...)

		if rsp.Links == nil || rsp.Links.IsLastPage() {
			return droplets, nil
		} else {
			page, err := rsp.Links.CurrentPage()
			if err != nil {
				return droplets, nil
			}
			opt.Page = page + 1
		}
	}
}

type tokenSource struct {
	AccessToken string
}

func (t *tokenSource) Token() (*oauth2.Token, error) {
	return &oauth2.Token{
		AccessToken: t.AccessToken,
	}, nil
}
