package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/t-ashula/go-narrow"
)

func main() {
	os.Exit(run())
}

func run() int {
	fmt.Fprintf(os.Stderr, "go-narrow client %s\n", narrow.Version)
	client := narrow.NewClient()
	params := narrow.NewSearchR18Params()
	params.SetLimit(1)
	params.AddOutputFields([]narrow.OutputField{narrow.OutputFieldNovelType})
	res, err := client.Search(context.Background(), params)
	if err != nil {
		fmt.Fprintf(os.Stderr, "search failed. %v", err)
		return 1
	}
	fmt.Fprintf(os.Stdout, "AllCount:%d\n", res.AllCount)
	for i, n := range res.NovelInfos {
		v, err := json.Marshal(&n)
		if err != nil {
			v = []byte(fmt.Sprintf("%+v", n))
		}
		fmt.Fprintf(os.Stdout, "%d : %s\n", i, v)
	}
	return 0
}
