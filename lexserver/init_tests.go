package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/stts-se/pronlex/lex"
)

func runInitTests(s *http.Server, port string) error {

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Fatal(fmt.Errorf("init_tests: couldn't start test server on port %s : %v", port, err))
		}
	}()

	log.Println("init_tests: waiting for server to start ...")

	time.Sleep(time.Second)

	log.Printf("init_tests: server up and running using port " + port)
	log.Println("init_tests: running tests")

	nErrs1, nTests1, err1 := testExampleURLs(port)
	nErrs2, nTests2, err2 := testURLsWithContent(port)

	var err error
	if err1 != nil && err2 != nil {
		err = fmt.Errorf("%v, %v", err1, err2)
	} else if err1 != nil {
		err = err1
	} else if err2 != nil {
		err = err2
	}
	if err != nil {
		return err
	}

	nTests := nTests1 + nTests2
	testString := "tests"
	if nTests == 1 {
		testString = "test"
	}
	if nErrs1 > 0 || nErrs2 > 0 {
		nErrs := nErrs1 + nErrs2
		errString := "errors"
		if nErrs == 1 {
			errString = "error"
		}
		log.Printf("init_tests: %d %s completed with %d %s!", nTests, testString, nErrs, errString)
		return fmt.Errorf("INIT TESTS FAILED")
	}

	log.Printf("init_tests: %d %s completed without errors", nTests, testString)
	return nil
}

func shortenURL(url string) string {
	limit := 50
	r := []rune(url)
	if len(r) > limit {
		return string(r[0:limit]) + " (...)"
	}
	return url
}

func testURLsWithContent(port string) (int, int, error) {

	log.Println("init_tests: testing some URLs with content")

	nFailed := 0
	nTests := 0

	lookupTests := map[string]string{
		"/lexicon/lookup?lexicons=demodb:demolex&wordlike=h%C3%A4st__": `[{"id":6,"lexRef":{"DBRef":"demodb","LexName":"demolex"},"strn":"hästar","language":"sv","partOfSpeech":"NN","morphology":"NEU IND PLU","wordParts":"hästar","lemma":{"id":4,"strn":"häst","reading":"","paradigm":""},"transcriptions":[{"id":9,"entryId":6,"strn":"\" h E . s t a r","language":"sv","sources":[]}],"status":{"id":6,"name":"demo","source":"auto","timestamp":"2017-08-18T09:37:51Z","current":true},"entryValidations":[],"preferred":false},{"id":7,"lexRef":{"DBRef":"demodb","LexName":"demolex"},"strn":"hästar","language":"sv","partOfSpeech":"NN","morphology":"NEU IND PLU","wordParts":"hästar","lemma":{"id":4,"strn":"häst","reading":"","paradigm":""},"transcriptions":[{"id":10,"entryId":7,"strn":"\" h { . s t a r","language":"sv","sources":[]}],"status":{"id":7,"name":"demo","source":"auto","timestamp":"2017-08-18T09:37:51Z","current":true},"entryValidations":[],"preferred":true}]`,

		"/lexicon/lookup?lexicons=demodb:demolex&wordpartsregexp=h%C3%A4st": `[{"id":5,"lexRef":{"DBRef":"demodb","LexName":"demolex"},"strn":"häst","language":"sv","partOfSpeech":"NN","morphology":"NEU IND SIN","wordParts":"häst","lemma":{"id":4,"strn":"häst","reading":"","paradigm":""},"transcriptions":[{"id":8,"entryId":5,"strn":"\" h E s t","language":"sv","sources":[]}],"status":{"id":5,"name":"demo","source":"auto","timestamp":"2017-08-18T10:02:32Z","current":true},"entryValidations":[],"preferred":true},{"id":6,"lexRef":{"DBRef":"demodb","LexName":"demolex"},"strn":"hästar","language":"sv","partOfSpeech":"NN","morphology":"NEU IND PLU","wordParts":"hästar","lemma":{"id":4,"strn":"häst","reading":"","paradigm":""},"transcriptions":[{"id":9,"entryId":6,"strn":"\" h E . s t a r","language":"sv","sources":[]}],"status":{"id":6,"name":"demo","source":"auto","timestamp":"2017-08-18T10:02:32Z","current":true},"entryValidations":[],"preferred":false},{"id":7,"lexRef":{"DBRef":"demodb","LexName":"demolex"},"strn":"hästar","language":"sv","partOfSpeech":"NN","morphology":"NEU IND PLU","wordParts":"hästar","lemma":{"id":4,"strn":"häst","reading":"","paradigm":""},"transcriptions":[{"id":10,"entryId":7,"strn":"\" h { . s t a r","language":"sv","sources":[]}],"status":{"id":7,"name":"demo","source":"auto","timestamp":"2017-08-18T10:02:32Z","current":true},"entryValidations":[],"preferred":true}]`,

		"/lexicon/lookup?lemmas=kex&lexicons=demodb:demolex": `[{"id":1,"lexRef":{"DBRef":"demodb","LexName":"demolex"},"strn":"kex","language":"sv","partOfSpeech":"NN","morphology":"NEU IND SIN","wordParts":"kex","lemma":{"id":1,"strn":"kex","reading":"","paradigm":""},"transcriptions":[ { "id":1,"entryId":1,"strn":"\" k e k s","language":"sv","sources":[]}, { "id":2,"entryId":1,"strn":"\" C e k s","language":"sv","sources":[]}],"status":{"id":1,"name":"demo","source":"auto","timestamp":"2017-08-18T09:57:30Z","current":true},"entryValidations":[],"preferred":true}, { "id":2,"lexRef":{"DBRef":"demodb","LexName":"demolex"},"strn":"kexet","language":"sv","partOfSpeech":"NN","morphology":"NEU DEF SIN","wordParts":"kexet","lemma":{"id":1,"strn":"kex","reading":"","paradigm":""},"transcriptions":[ { "id":3,"entryId":2,"strn":"\" k e k . s @ t","language":"sv","sources":[]}, { "id":4,"entryId":2,"strn":"\" C e k . s @ t","language":"sv","sources":[]}],"status":{"id":2,"name":"demo","source":"auto","timestamp":"2017-08-18T09:57:30Z","current":true},"entryValidations":[],"preferred":true}]`,

		"/lexicon/lookup?lexicons=demodb:demolex&words=dom&transcriptionlike=%25o:%25&pp=yes": `[
 { "id": 9, "lexRef": { "DBRef": "demodb", "LexName": "demolex" }, "strn": "dom", "language": "sv", "partOfSpeech": "NN", "morphology": "UTR IND SIN", "wordParts": "dom", "lemma": { "id": 5, "strn": "dom", "reading": "", "paradigm": "" }, "transcriptions": [ { "id": 12, "entryId": 9, "strn": "\" d o: m", "language": "sv", "sources": [] } ], "status": { "id": 9, "name": "demo", "source": "auto", "timestamp": "2017-08-18T10:04:44Z", "current": true }, "entryValidations": [], "preferred": true }]`,
	}

	jsonMapTests := map[string]string{
		"/symbolset/list": `{"symbol_set_names":["ar_sampa_mary","ar_ws-sampa","en-us_cmu","en-us_sampa_mary","en-us_ws-sampa","nb-no_nst-xsampa","nb-no_ws-sampa","sv-se_nst-xsampa","sv-se_sampa_mary","sv-se_ws-sampa"]}`,
		"/mapper/map/sv-se_ws-sampa/sv-se_sampa_mary/%22%22%20p%20O%20j%20.%20k%20@": `{"From":"sv-se_ws-sampa","To":"sv-se_sampa_mary","Input":"\"\" p O j . k @","Result":"\" p O j - k @"}`,
		"/mapper/map/sv-se_sampa_mary/sv-se_ws-sampa/%22%20p%20O%20j%20-%20k%20@":    `{"From":"sv-se_sampa_mary","To":"sv-se_ws-sampa","Input":"\" p O j - k @","Result":"\"\" p O j . k @"}`,
		"/validation/list": `{"validator_names":["en-us_ws-sampa","nb-no_ws-sampa","sv-se_ws-sampa"]}`,
		`/validation/validateentry?symbolsetname=en-us_ws-sampa&entry={%22id%22:1703348,%22lexiconId%22:3,%22strn%22:%22barn%22,%22language%22:%22en-us%22,%22partOfSpeech%22:%22%22,%22wordParts%22:%22%22,%22lemma%22:{%22id%22:0,%22strn%22:%22%22,%22reading%22:%22%22,%22paradigm%22:%22%22},%22transcriptions%22:[{%22id%22:1717337,%22entryId%22:1703348,%22strn%22:%22\%22%20b%20A%20r%20n%22,%22language%22:%22%22,%22sources%22:[]}],%22status%22:{%22id%22:1703348,%22name%22:%22imported%22,%22source%22:%22cmu%22,%22timestamp%22:%222016-09-06T13:16:07Z%22,%22current%22:true},%22entryValidations%22:[]}`: `{"id":1703348,"lexRef":{"DBRef":"","LexName":""},"strn":"barn","language":"en-us","partOfSpeech":"","morphology":"","wordParts":"","lemma":{"id":0,"strn":"","reading":"","paradigm":""},"transcriptions":[{"id":1717337,"entryId":1703348,"strn":"\" b A r n","language":"","sources":[]}],"status":{"id":1703348,"name":"imported","source":"cmu","timestamp":"2016-09-06T13:16:07Z","current":true},"entryValidations":[{"id":0,"level":"Format","ruleName":"primary_stress","Message":"Each trans should have one primary stress. Found: /\" b A r n/","timestamp":""},{"id":0,"level":"Fatal","ruleName":"SymbolSet","Message":"Invalid transcription symbol '\"' in /\" b A r n/","timestamp":""}],"preferred":false}`,
	}

	jsonListTestsMustContain := map[string][]string{
		"/admin/list_dbs": []string{"demodb"},
		"/mapper/list":    []string{"sv-se_ws-sampa - sv-se_sampa_mary", "sv-se_sampa_mary - sv-se_ws-sampa"},
	}

	jsonBoolTests := map[string]bool{
		"/validation/has_validator/sv-se_ws-sampa": true,
		"/validation/has_validator/ar_ws-sampa":    false,
	}

	log.Printf("init_tests: testing entry lookup: %d", len(lookupTests))
	for url, expect := range lookupTests {
		nTests = nTests + 1
		ok, err := lookupTest(port, url, expect)
		if !ok {
			nFailed = nFailed + 1
		}
		if err != nil {
			return nFailed, nTests, err
		}
	}

	log.Printf("init_tests: testing json map results: %d", len(jsonMapTests))
	for url, expect := range jsonMapTests {
		nTests = nTests + 1
		ok, err := jsonMapTest(port, url, expect)
		if !ok {
			nFailed = nFailed + 1
		}
		if err != nil {
			return nFailed, nTests, err
		}
	}

	log.Printf("init_tests: testing json list results: %d", len(jsonListTestsMustContain))
	for url, expect := range jsonListTestsMustContain {
		nTests = nTests + 1
		ok, err := jsonListTestMustContain(port, url, expect)
		if !ok {
			nFailed = nFailed + 1
		}
		if err != nil {
			return nFailed, nTests, err
		}
	}

	log.Printf("init_tests: testing boolean results: %d", len(jsonBoolTests))
	for url, expect := range jsonBoolTests {
		nTests = nTests + 1
		ok, err := jsonTestBool(port, url, expect)
		if !ok {
			nFailed = nFailed + 1
		}
		if err != nil {
			return nFailed, nTests, err
		}
	}

	return nFailed, nTests, nil
}

func jsonMapTest(port string, url string, expect string) (bool, error) {
	url = "http://localhost" + port + url
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("** FAILED TEST ** for %s : couldn't retreive URL : %v\n", url, err)
		return false, nil
	} else {
		log.Printf("init_tests: jsonMap %s", url)

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("** FAILED TEST ** for %s : expected response code 200, found %d\n", url, resp.StatusCode)
			return false, nil
		}

		got, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false, fmt.Errorf("couldn't read response body : %v", err)
		}

		var gotJ map[string]interface{}
		err = json.Unmarshal([]byte(got), &gotJ)
		if err != nil {
			return false, fmt.Errorf("couldn't convert response to json : %v", err)
		}

		var expJ map[string]interface{}
		err = json.Unmarshal([]byte(expect), &expJ)
		if err != nil {
			return false, fmt.Errorf("couldn't convert expected result to json : %v", err)
		}

		if !reflect.DeepEqual(gotJ, expJ) {
			fmt.Printf("** FAILED TEST ** for %s :\n >> EXPECTED RESPONSE:\n%s\n >> FOUND:\n%s\n", url, expect, string(got))
			return false, nil
		}
	}
	return true, nil
}

func jsonTestBool(port string, url string, expect bool) (bool, error) {
	url = "http://localhost" + port + url
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("** FAILED TEST ** for %s : couldn't retreive URL : %v\n", url, err)
		return false, nil
	} else {
		log.Printf("init_tests: jsonMap %s", url)

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("** FAILED TEST ** for %s : expected response code 200, found %d\n", url, resp.StatusCode)
			return false, nil
		}

		got, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false, fmt.Errorf("couldn't read response body : %v", err)
		}

		var gotJ bool
		err = json.Unmarshal([]byte(got), &gotJ)
		if err != nil {
			return false, fmt.Errorf("couldn't convert response to json : %v", err)
		}
		if gotJ != expect {
			fmt.Printf("** FAILED TEST ** for %s :\n >> EXPECTED RESPONSE: %v\n >> FOUND: %v\n", url, expect, gotJ)
			return false, nil
		}
	}
	return true, nil
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func jsonListTestMustContain(port string, url string, expect []string) (bool, error) {
	url = "http://localhost" + port + url
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("** FAILED TEST ** for %s : couldn't retreive URL : %v\n", url, err)
		return false, nil
	} else {
		log.Printf("init_tests: jsonList %s", url)

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("** FAILED TEST ** for %s : expected response code 200, found %d\n", url, resp.StatusCode)
			return false, nil
		}

		got, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false, fmt.Errorf("couldn't read response body : %v", err)
		}

		var gotJ []string
		err = json.Unmarshal([]byte(got), &gotJ)
		if err != nil {
			return false, fmt.Errorf("couldn't convert response to json : %v", err)
		}

		ok := true
		for _, exp := range expect {
			if !contains(gotJ, exp) {
				ok = false
			}
		}
		if !ok {
			fmt.Printf("** FAILED TEST ** for %s :\n >> EXPECTED RESPONSE TO CONTAIN ALL OF:\n%s\n >> FOUND:\n%s\n", url, expect, string(got))
			return false, nil
		}

	}
	return true, nil
}

func lookupTest(port string, url string, expect string) (bool, error) {
	url = "http://localhost" + port + url
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		fmt.Printf("** FAILED TEST ** for %s : couldn't retreive URL : %v\n", url, err)
		return false, nil
	} else {
		log.Printf("init_tests: lookup/entry %s", url)

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("** FAILED TEST ** for %s : expected response code 200, found %d\n", url, resp.StatusCode)
			return false, nil
		}

		got, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return false, fmt.Errorf("couldn't read response body : %v", err)
		}
		var gotEs []lex.Entry
		err = json.Unmarshal(got, &gotEs)
		if err != nil {
			return false, fmt.Errorf("couldn't parse json : %v", err)
		}
		var expEs []lex.Entry
		err = json.Unmarshal([]byte(expect), &expEs)
		if err != nil {
			return false, fmt.Errorf("couldn't parse expect string : %v", err)
		}
		for i, e := range gotEs {
			e.EntryStatus.Timestamp = ""
			gotEs[i] = e
		}
		for i, e := range expEs {
			e.EntryStatus.Timestamp = ""
			expEs[i] = e
		}

		if !reflect.DeepEqual(gotEs, expEs) {
			fmt.Printf("** FAILED TEST ** for %s :\n >> EXPECTED RESPONSE:\n%s\n >> FOUND:\n%s\n", url, expect, string(got))
			return false, nil
		}
	}
	return true, nil
}

func testExampleURLs(port string) (int, int, error) {

	log.Println("init_tests: testing response codes for built-in example URLs")

	nFailed := 0
	nTests := 0

	resp, err := http.Get("http://localhost" + port + "/meta/examples")
	defer resp.Body.Close()
	if err != nil {
		return nFailed, nTests, fmt.Errorf("couldn't retrieve server's url examples : %v", err)
	}
	js, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nFailed, nTests, fmt.Errorf("couldn't retrieve server's url examples : %v", err)
	}

	var res []JSONURLExample
	err = json.Unmarshal([]byte(js), &res)
	if err != nil {
		return nFailed, nTests, fmt.Errorf("couldn't unmarshal json : %v", err)
	}

	for _, example := range res {
		nTests = nTests + 1
		url := "http://localhost" + port + urlEnc(example.URL)
		resp, err = http.Get(url)
		defer resp.Body.Close()
		if err != nil {
			fmt.Printf("** FAILED TEST ** for %s : couldn't retreive URL : %v\n", url, err)
			nFailed = nFailed + 1
		} else {
			log.Printf("init_tests: %s => %s : %s", example.Template, shortenURL(example.URL), resp.Status)

			if resp.StatusCode != http.StatusOK {
				fmt.Printf("** FAILED TEST ** for %s : expected response code 200, found %d\n", url, resp.StatusCode)
				nFailed = nFailed + 1
			} else {
				got, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					return nFailed, nTests, fmt.Errorf("couldn't read response body : %v", err)
				}

				if strings.TrimSpace(string(got)) == "" {
					fmt.Printf("** FAILED TEST ** for %s : expected non-empty response\n", url, resp)
					nFailed = nFailed + 1
				}
			}
		}
	}

	return nFailed, nTests, nil
}
