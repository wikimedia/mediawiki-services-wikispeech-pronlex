package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	usage := `demotest usage:
		$ go run demotest.go <PORT>`

	if len(os.Args) != 2 {
		fmt.Println(usage)
		os.Exit(1)
	}
	port := os.Args[1]
	log.Println("demotest: running tests...")
	err := runTests(port)
	if err != nil {
		log.Printf("demotest: tests failed : %v", err)
		os.Exit(1)
	}
	log.Println("demotest: all tests completed")
}

type JSONURLExample struct {
	Template string `json:"template"`
	URL      string `json:"url"`
}

// TODO: Neat URL encoding...
func urlEnc(url string) string {
	return strings.Replace(strings.Replace(strings.Replace(url, " ", "%20", -1), "\n", "", -1), `"`, "%22", -1)
}

func shortenURL(url string) string {
	limit := 50
	r := []rune(url)
	if len(r) > limit {
		return string(r[0:limit]) + " (...)"
	}
	return url
}

func runTests(port string) error {

	resp, err := http.Get("http://localhost:" + port + "/meta/examples")
	defer resp.Body.Close()
	if err != nil {
		return fmt.Errorf("couldn't retrieve server's url examples : %v", err)
	}
	js, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("couldn't retrieve server's url examples : %v", err)
		os.Exit(1)
	}

	var res []JSONURLExample
	err = json.Unmarshal([]byte(js), &res)
	if err != nil {
		return fmt.Errorf("couldn't unmarshal json : %v", err)
	}

	for _, example := range res {
		url := "http://localhost:" + port + urlEnc(example.URL)
		resp, err = http.Get(url)
		defer resp.Body.Close()
		if err != nil {
			return fmt.Errorf("couldn't get URL : %v", err)
		}
		log.Printf("demotest: %s => %s : %s", example.Template, shortenURL(example.URL), resp.Status)

		if resp.StatusCode != http.StatusOK {
			return fmt.Errorf("expected response status 200, found %s for url %s", resp.Status, url)
			os.Exit(1)
		}

		// TODO: Read body?
		// r, err := ioutil.ReadAll(resp.Body)
		// if err != nil {
		// 	return fmt.Errorf("couldn't read response body : %v", err)
		// }
		//fmt.Println(string(r))
	}

	return nil
}
