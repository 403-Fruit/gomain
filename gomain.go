package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/tj/go-spin"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
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
	endpoint           = "https://domainr.com/api/json/search?q=%s&client_id=gomain"
	out      io.Writer = os.Stdout
	box                = spin.Box1
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
		fmt.Fprintf(os.Stderr, "Domain required.\n")
		os.Exit(1)
	}

	wg := sync.WaitGroup{}
	ex := make(chan struct{})

	go func() {
		wg.Add(1)
		tick(50*time.Millisecond, ex)
		wg.Done()
	}()

	req, err := http.Get(fmt.Sprintf(endpoint, argv[0]))

	close(ex)
	wg.Wait()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Server could not be reached.\n")
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid response from server.\n")
		os.Exit(1)
	}

	var query query

	json.Unmarshal(body, &query)

	for _, dom := range query.Results {
		switch dom.Availability {
		case "available":
			fmt.Fprintf(out, "%sA %s%s\n", green, normal, dom.Domain)
		case "maybe", "unknown":
			fmt.Fprintf(out, "%sM %s%s\n", yellow, normal, dom.Domain)
		default:
			fmt.Fprintf(out, "%sU %s%s\n", gray, normal, dom.Domain)
		}
	}
}

func tick(d time.Duration, ex chan struct{}) {
	t := time.NewTicker(d)
	s := spin.New()
	s.Set(box)

	for {
		select {
		case <-t.C:
			fmt.Fprintf(out, "\r%s", s.Next())
		case <-ex:
			t.Stop()
			fmt.Printf("\r")
			return
		}
	}
}
