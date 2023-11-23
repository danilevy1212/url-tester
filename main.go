package main

import (
	"bufio"
	"fmt"
	"github.com/alecthomas/kingpin/v2"
	"log"
	urllib "net/url"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"url-tester/lib/url"
	"url-tester/lib/worker"
)

var (
	app  = kingpin.New("url-tester", "Reports not available URLs by sending HEAD requests to them.")
	file = app.Flag("file", "Test all urls contained in a file, separated by a new line").Short('f').String()
	link = app.Arg("url", "URL to send a HEAD request to").String()
)

// Must be called after parsing arguments
func getUrls() ([]string, error) {
	var urls []string

	if *file == "" && *link == "" {
		app.Fatalf("No urls passed, see --help")
	}

	if *link != "" {
		argURL, err := url.SanitizeURL(*link)
		if err != nil {
			urls = append(urls, argURL)
		} else {
			return nil, fmt.Errorf("invalid url %s: %s", *link, err)
		}
	}

	if *file == "" {
		return urls, nil
	}

	path, err := filepath.Abs(*file)
	if err != nil {
		return nil, fmt.Errorf("invalid filename %s: %s", *file, err)
	}

	fileinfo, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("file %s does not exist: %s", path, err)
	}
	if fileinfo.IsDir() {
		return nil, fmt.Errorf("given path %s is a directory", path)
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %s", err)
	}
	defer file.Close()

	urls, err = parseURLFile(file, urls)
	if err != nil {
		return nil, fmt.Errorf("error parsing provided file: %s", err)
	}

	return urls, nil
}

func parseURLFile(fileHandler *os.File, result []string) ([]string, error) {
	scanner := bufio.NewScanner(fileHandler)
	originalWriter := log.Writer()
	defer log.SetOutput(originalWriter)

	log.SetOutput(os.Stderr)

	for scanner.Scan() {
		rawURL := scanner.Text()
		url, err := url.SanitizeURL(rawURL)

		if err == nil {
			result = append(result, url)
		} else {
			log.Printf("WARNING: error ocurred during parsing of url `%s`, ignoring: %s", rawURL, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %s", err)
	}

	return result, nil
}

func main() {
	kingpin.MustParse(app.Parse(os.Args[1:]))

	urls, err := getUrls()
	if err != nil {
		log.Fatal("Error parsing URL file:", err)
	}

	var wg sync.WaitGroup
	wg.Add(len(urls))
	numWorkers := runtime.NumCPU() * 2
	results := make(chan worker.HeadResult, len(urls))
	url := make(chan string, len(urls))

	// Setup workers
	for i := 0; i < numWorkers; i++ {
		go worker.Worker(url, results, &wg)
	}

	// Send urls to them
	for _, u := range urls {
		url <- u
	}
	close(url)

	// When all urls are checked, it is safe to close.
	wg.Wait()
	close(results)

	for r := range results {
		u, _ := urllib.QueryUnescape(r.URL)

		fmt.Printf("Status '%d': '%s'\n", r.StatusCode, u)
	}
}
