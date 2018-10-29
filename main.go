package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	httpstat "github.com/tcnksm/go-httpstat"
)

func main() {
	// Open aws.json
	jsonFile, err := os.Open("./endpoints/aws.json")

	// check err
	if err != nil {
		fmt.Println(err)
		os.Exit(42)
	}

	fmt.Println("Successfully Opened aws.json")
	// defer closing
	defer jsonFile.Close()

	var endpoints map[string]string

	// Read our file
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &endpoints)

	for k, v := range endpoints {
		// fmt.Println(k)
		// fmt.Println(v)
		pingEndpoint(k, v)
	}

}

func pingEndpoint(name string, endpoint string) {
	// Create a new HTTP request
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Create a httpstat powered context
	var result httpstat.Result
	ctx := httpstat.WithHTTPStat(req.Context(), &result)
	req = req.WithContext(ctx)

	// Send request by default HTTP client
	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	if _, err := io.Copy(ioutil.Discard, res.Body); err != nil {
		log.Fatal(err)
	}
	res.Body.Close()

	// // Show the results
	// log.Printf("DNS lookup: %d ms", int(result.DNSLookup/time.Millisecond))
	log.Printf("%s - %s : %d ms", name, endpoint, int(result.TCPConnection/time.Millisecond))
	// log.Printf("TLS handshake: %d ms", int(result.TLSHandshake/time.Millisecond))
	// log.Printf("Server processing: %d ms", int(result.ServerProcessing/time.Millisecond))
	// log.Printf("Content transfer: %d ms", int(result.ContentTransfer(time.Now())/time.Millisecond))
	// log.Printf("result.Total %d ms \n\n", int(result.Total(time.Now())/time.Millisecond))

}
