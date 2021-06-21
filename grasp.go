package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

var parallelFactor = flag.Int("parallel", 10, "")

// normalizeAddresses adds http:// prefix whenever it is necessary.
func normalizeAddresses(addresses []string) []string {
	result := make([]string, 0)
	for _, addr := range addresses {
		normalizedAddress := addr
		if !strings.HasPrefix(addr, "http://") &&
			!strings.HasPrefix(addr, "https://") {
			normalizedAddress = fmt.Sprintf("http://%s", addr)
		}
		result = append(result, normalizedAddress)
	}
	return result
}

// fetcher is the worker, it fetches the page and
// produces the md5 sum to the printChannel channel
func fetcher(
  inputChannel chan string,
  printChannel chan string,
  wg *sync.WaitGroup,
  client http.Client,
) {
	defer wg.Done()
	for addr := range inputChannel {
		// Instead of having an error channel, I decided to send
		// any error messages out using the same print channel
		resp, err := client.Get(addr)
		if err != nil {
			printChannel <- fmt.Sprintf("%s could not perform request", addr)
			continue
		}
		raw, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			printChannel <- fmt.Sprintf("%s could not read response body", addr)
		}
		printChannel <- fmt.Sprintf("%s %x", addr, md5.Sum(raw))
		_ = resp.Body.Close()
	}
}

func run(client http.Client, addresses []string) error {
	var wgFetcher sync.WaitGroup
	var wgPrinter sync.WaitGroup

	inputChannel := make(chan string)
	printChannel := make(chan string)

	// starting fetcher workers
	for workers := 0; workers < *parallelFactor; workers++ {
		wgFetcher.Add(1)
		go fetcher(inputChannel, printChannel, &wgFetcher, client)
	}

	// starting printer worker
	wgPrinter.Add(1)
	go func() {
		defer wgPrinter.Done()
		for msg := range printChannel {
			fmt.Println(msg)
		}
	}()
	for _, addr := range addresses {
		inputChannel <- addr
	}
	close(inputChannel)
	wgFetcher.Wait()
	close(printChannel)
	wgPrinter.Wait()

	return nil
}

func main() {
	flag.Parse()
	if err := run(http.Client{}, normalizeAddresses(flag.Args())); err != nil {
		fmt.Println(err)
	}
}

