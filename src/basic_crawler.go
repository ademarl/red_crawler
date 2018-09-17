
package main

import "time"


//===========================================================================
// Defines the behaviour of the crawler. Channels are used as semaphores to
// speedup the execution while waiting for http replies
//===========================================================================


// Wait time between requests to the target website in milliseconds. Might need calibration
const WAIT = 100

// Green light for semaphores used in the crawl_shares function
type green struct {}


//===========================================================================


// Extracts the targeted information from all URLs by crawling through the links with parallel requests
func Crawl_shares(links []string, papers []string, n int) []Paper {

	var shares = make([]Paper, n)

	// Parallel crawling with channels as semaphores
	semaphore := make(chan green, n);

	for i := range papers {
		go Share_parser(shares, links[i], papers[i], i, semaphore)
		// Server is unhappy with too many requests, slow it down
		time.Sleep(WAIT*time.Millisecond)
	}

	// Wait until all semaphores are green
	for i := 0; i < n; i++ { <-semaphore }

	return shares
}


//===========================================================================
