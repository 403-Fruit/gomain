package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/tj/go-spin"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const usage = `
  Usage: gomain <domain>

  Examples:

  $ gomain daryl
  $ gomain daryl.im
  $ gomain foo.com
`

const (
	gray   = "\033[0;37m"
	green  = "\033[0;32m"
	yellow = "\033[0;33m"
	normal = "\033[0;00m"
)

var (
	endpoint = "https://domainr.com/api/json/search?q=%s&client_id=gomain"
)

type query struct {
	Results []result
	Query   string
}

type result struct {
	Availability string
	Domain       string
}

func main() {
	flag.Parse()

	flag.Usage = func() {
		fmt.Println(usage)
		os.Exit(0)
	}

	argv := flag.Args()

	if len(argv) < 1 {
		fmt.Printf("Domain required.\n")
		return
	}

	tick := time.NewTicker(50 * time.Millisecond)

	go tock(tick)
	req, err := http.Get(fmt.Sprintf(endpoint, argv[0]))
	fmt.Printf("\n")
	tick.Stop()

	if err != nil {
		fmt.Printf("Server could not be reached.\n")
		return
	}

	var query query

	body, _ := ioutil.ReadAll(req.Body)

	json.Unmarshal(body, &query)

	for _, dom := range query.Results {
		switch dom.Availability {
		case "available":
			fmt.Printf("%sA %s%s\n", green, normal, dom.Domain)
		case "maybe", "unknown":
			fmt.Printf("%sM %s%s\n", yellow, normal, dom.Domain)
		default:
			fmt.Printf("%sU %s%s\n", gray, normal, dom.Domain)
		}
	}
}

func tock(t *time.Ticker) {
	s := spin.New()
	s.Set(spin.Box1)

	for _ = range t.C {
		fmt.Printf("\r%s", s.Next())
	}
}
