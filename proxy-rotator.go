package main

import (
	"bufio"
	"log"
	"math/rand"
	"os"
	"time"
)

func ProxyRotator(useProxy bool) string {
	if !useProxy {
		return ""
	}

	// Create a new source of random numbers seeded with the current time
	source := rand.NewSource(time.Now().UnixNano())

	// Create a new random number generator using the source
	random := rand.New(source)

	path, _ := os.Getwd()

	// Open the file
	file, err := os.Open(path + "./proxies.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Read the file line by line and store the proxies in a slice
	var proxies []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		proxies = append(proxies, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Select a random proxy from the slice
	randomIndex := random.Intn(len(proxies))
	return proxies[randomIndex]
}
