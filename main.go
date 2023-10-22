package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

func parseArgs() (string, error) {
	flag.Parse()

	filename := flag.Arg(0)

	if filename == "" {
		return "", fmt.Errorf("no filepath specified")
	}

	path, err := filepath.Abs(filename)
	if err != nil {
		return "", fmt.Errorf("invalid filename %s: %s", filename, err)
	}

	_, err = os.Stat(path)
	if err != nil {
		return "", fmt.Errorf("file %s does not exist: %s", path, err)
	}

	return path, nil
}

func parseURLFile(fileHandler *os.File) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(fileHandler)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %s", err)
	}

	return lines, nil
}

func searchWithHead(url string, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Head(url)
	if err != nil {
		log.Printf("Error: %s\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("`%s` was not found: status %d\n", url, resp.StatusCode)
	}
}

func main() {
	path, err := parseArgs()
	if err != nil {
		log.Fatal("Error:", err)
	}

	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Error opening file:", err)
	}
	defer file.Close()

	var wg sync.WaitGroup
	urls, err := parseURLFile(file)
	if err != nil {
		log.Fatal("Error parsing URL file:", err)
	}

	wg.Add(len(urls))

	for _, u := range urls {
		go searchWithHead(u, &wg)
	}

	wg.Wait()
}
