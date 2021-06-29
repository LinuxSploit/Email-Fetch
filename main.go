package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"sync"
	"time"
)

func getemails(url string) {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("connection error: ", url)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("unable to read response body: ", url)
		return
	}
	bodystr := string(body)

	//      Email extraction regular expresion:
	//      [a-z0-9\.\-+_]+@[a-z0-9\.\-+_]+\.[a-z]+

	var validID = regexp.MustCompile(`[a-z0-9\.\-+_]+@[a-z0-9\.\-+_]+\.[a-z]+`)
	emails := validID.FindAllString(bodystr, -1)

	if len(emails) == 0 {
		return
	}

	for x := range emails {
		fmt.Println(emails[x])
	}

	return
}
func currrnt(d []string) {
	for _, x := range d {
		wg.Add(1)
		go getemails(x)
	}
	wg.Wait()
}

var wg sync.WaitGroup

func main() {
	started := time.Now()
	var domains []string

	/////////////FLAG
	if len(os.Args) < 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Printf("usage: %v <concurrency>\n", os.Args[0])
		os.Exit(1)
	}

	var concurrnt int
	_, e := fmt.Sscan(os.Args[1], &concurrnt)
	if e != nil {
		fmt.Printf("usage: %v <concurrency>\n", os.Args[0])
		log.Fatal("invalid <concurrency> value")
	}

	/////////////STDIN
	scanoye := bufio.NewScanner(os.Stdin)
	for scanoye.Scan() {
		domains = append(domains, scanoye.Text())
	}

	/////////////CONCURRENCY-VALUE-CHECK
	if concurrnt < 2 || concurrnt > len(domains) {
		concurrnt = 1
	}

	/////////////
	for x := 0; x < len(domains); x = x + concurrnt {
		if x == ((len(domains) / concurrnt) * concurrnt) {
			currrnt(domains[x:len(domains)])
		} else {
			currrnt(domains[x : x+concurrnt])
		}
	}
	fmt.Println(time.Since(started))
	/////////////
}
