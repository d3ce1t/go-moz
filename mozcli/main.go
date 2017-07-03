package main

import (
	"flag"
	"fmt"
	"os"
	"peeple/seo-tools/mozapi"
)

type arrayFlags []string

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func (i *arrayFlags) String() string {
	return fmt.Sprintf("%v", *i)
}

var urls arrayFlags

func main() {

	flag.Var(&urls, "u", "URL to retrieve information for.")
	flag.Parse()

	if flag.NArg() > 0 || len(urls) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	moz := mozapi.New("ACCESS_ID", "SECRET_KEY")
	columns := mozapi.Title | mozapi.CanonicalURL | mozapi.ExternalEquityLinks |
		mozapi.Links | mozapi.MozRankForURL | mozapi.MozRankForSubdomain |
		mozapi.HTTPStatusCode | mozapi.PageAuthority | mozapi.DomainAuthority |
		mozapi.TimeLastCrawled

	metrics, err := moz.MetricsForURLBatch(urls, columns)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(-1)
	}

	for _, metric := range metrics {
		fmt.Println(metric)
	}
}
