package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
)

func getemails(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("URL:", url)
		fmt.Println("    unable to connect")
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("URL:", url)
		fmt.Println("    unable to read HTML-BODY")
		return
	}
	bodystr := string(body)

	//      Email extraction regular expresion:
	//      [a-z0-9\.\-+_]+@[a-z0-9\.\-+_]+\.[a-z]+

	var validID = regexp.MustCompile(`[a-z0-9\.\-+_]+@[a-z0-9\.\-+_]+\.[a-z]+`)
	emails := validID.FindAllString(bodystr, -1)

	fmt.Printf("URL: %v  \n", url)
	if len(emails) == 0 {
		fmt.Println("    no email found!")
		return
	}

	for x := range emails {
		fmt.Println("    " + emails[x])
	}

}
func main() {

	scanoye := bufio.NewScanner(os.Stdin)
	for scanoye.Scan() {
		domain := scanoye.Text()
		getemails(domain)
	}

}
