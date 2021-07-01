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
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
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

var wg sync.WaitGroup

type domains []string

func (d *domains) append(url string) {
	*d = append(*d, url)
}
func (d *domains) empty() {
	*d = []string{}
}
func bestconcurrencyvalue() {
	//under development

}
func main() {
	var domains domains
	var v int
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
		domains.append(scanoye.Text())
	}
	len := len(domains)
	/////////////CONCURRENCY-VALUE-CHECK
	if concurrnt < 2 || concurrnt > len {
		concurrnt = 1
	}
	started := time.Now()

	/////////////
	for x := 0; x < len; x = x + concurrnt {
		v = ((len / concurrnt) * concurrnt)
		if x == v {
			///
			for _, x := range domains[x:len] {
				wg.Add(1)
				go getemails(x)
			}
			wg.Wait()
			///
		} else if x != v {
			///
			for _, x := range domains[x : x+concurrnt] {
				wg.Add(1)
				go getemails(x)
			}
			wg.Wait()
			///
		}
	}
	fmt.Println(" time: ", time.Since(started))
	/////////////
}
