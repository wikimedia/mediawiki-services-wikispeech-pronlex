package main

import (
	"fmt"
	"log"
	"net"
)

func serverInitTests() error {
	port := "8799"
	s, err := createServer(port)
	if err != nil {
		return err
	}

	_, err = net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return err
	}
	s.ListenAndServe()

	err = runTests(port)
	if err != nil {
		return err
	}

	s.Close()

	log.Println("init_tests: done")

	return nil
}

func runTests(port string) error {
	log.Println("init_tests: running tests ... (work in progress)")
	for _, subRouter := range subRouters {
		for _, handler := range subRouter.handlers {
			for _, example := range handler.examples {
				url := "http://localhost:" + port + subRouter.root + example
				template := subRouter.root + handler.url
				log.Println("init_tests: " + template + " => " + url)
				// resp, err := http.Get(url)
				// if err != nil {
				// 	return err
				// }
				// resp.Body.Close()
				// body, err := ioutil.ReadAll(resp.Body)
				// if err != nil {
				// 	return err
				// }
				// fmt.Println(body)
			}
		}
	}
	return nil
}
