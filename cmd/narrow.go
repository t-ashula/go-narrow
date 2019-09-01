package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/t-ashula/go-narrow"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Version = narrow.Version
	app.Commands = []cli.Command{
		searchCommand(),
		fetchCommand(),
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func fetchCommand() cli.Command {
	return cli.Command{
		Name:  "fetch",
		Usage: "Fetch from syosetu.com group",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "site",
				Value: "narou",
				Usage: "fetch from `SITE` {narou, noc(Nocturne), mid(midnight), ml(moonlight), mlbl(moonlight bl)}",
			},
			cli.StringFlag{
				Required: true,
				Name:     "ncode",
				Usage:    "fetch `NCODE`",
			},
			cli.BoolFlag{
				Name: "over18",
			},
		},
		Action: func(c *cli.Context) error {
			params := makeFetchParams(c)
			client := narrow.NewClient()
			res, err := client.Fetch(context.Background(), params)
			if err != nil {
				return err
			}
			fmt.Printf("FetchResult:%+v\n", res)
			return nil
		},
	}
}

func searchCommand() cli.Command {
	return cli.Command{
		Name:  "search",
		Usage: "Search from syosetu.com group",
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "site",
				Value: "narou",
				Usage: "search from `SITE` {narou, noc(Nocturne), mid(midnight), ml(moonlight), mlbl(moonlight bl)}",
			},
			cli.IntFlag{Name: "limit", Value: 20, Usage: "max number of output"},
			cli.IntFlag{Name: "start", Value: 1, Usage: "start from"},
			cli.StringFlag{
				Name:  "order",
				Value: "new",
				Usage: "order by `ITEM` {new, fav, review, hyoka, hyokaasc, impression, hyokacnt, hyokacntasc, weekly, lengthdesc, lengthasc, ncodedesc, old}",
			},
			cli.StringSliceFlag{
				Name:  "keywords",
				Usage: "search keywords",
			},
		},
		Action: func(c *cli.Context) error {
			site := c.String("site")
			if !isKnownSite(site) {
				return fmt.Errorf("unknown site `%s` specified", site)
			}
			params := makeSearchParams(c)
			client := narrow.NewClient()
			res, err := client.Search(context.Background(), params)
			if err != nil {
				return err
			}
			fmt.Fprintf(os.Stdout, "AllCount:%d\n", res.AllCount)
			for i, n := range res.NovelInfos {
				v, err := json.Marshal(&n)
				if err != nil {
					v = []byte(fmt.Sprintf("%+v", n))
				}
				fmt.Fprintf(os.Stdout, "%d : %s\n", i, v)
			}

			return nil
		},
	}
}

func isKnownSite(site string) bool {
	return site == "noc" || site == "mid" || site == "ml" || site == "mlbl" || site == "narou" || site == ""
}

func makeSearchParams(c *cli.Context) narrow.Params {
	params := narrow.NewSearchParams()
	if c.IsSet("limit") {
		params.SetLimit(c.Int("limit"))
	}
	if c.IsSet("start") {
		params.SetStart(c.Int("start"))
	}
	if c.IsSet("keywords") {
		params.AddWords(c.StringSlice("keywords"))
	}
	// TODO: more flags

	site := c.String("site")
	if site == "noc" || site == "mid" || site == "ml" || site == "mlbl" {
		r18 := narrow.NewSearchR18Params()
		r18.SearchParams = *params
		r18.AddNocGenres([]narrow.NocGenre{nocgenre(site)})
		// TODO: NotNocGenere,
		return r18
	}

	return params
}

func nocgenre(site string) narrow.NocGenre {
	switch site {
	case "noc":
		return narrow.NocGenreNocturne
	case "mid":
		return narrow.NocGenreMidnight
	case "ml":
		return narrow.NocGenreMoonlightWomen
	case "mlbl":
		return narrow.NocGenreMoonlightBL
	default:
		return narrow.NocGenreAll
	}
}

func makeFetchParams(c *cli.Context) *narrow.FetchParams {
	params := narrow.NewFetchParams()
	site := c.String("site")
	switch site {
	case "noc":
		params.Site = narrow.FetchSiteNocturne
	case "mid":
		params.Site = narrow.FetchSiteMidNight
	case "ml":
		params.Site = narrow.FetchSiteMoonLight
	case "mlbl":
		params.Site = narrow.FetchSiteMoonLight
	default:
		params.Site = narrow.FetchSiteNarou
	}
	params.NCode = c.String("ncode")
	params.AllowOver18 = c.Bool("over18")
	return params
}
