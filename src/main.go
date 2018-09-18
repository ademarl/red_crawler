
package main

import (
	"fmt"
	"red_crawler/src/red_crawler"
)


//===========================================================================
// Main method - uses sub package 'red_crawler'
// The application is a web crawler that searches for information of stock
// market shares of Bovespa on Fundamentus webpage
// After parsing and crawling, the top ten shares by market value are
// selected and inserted in a MySQL database
//
// Additionally to the specification, this application crawls in parallel
//===========================================================================


// Main function
func main() {

	var shares []red_crawler.Paper

	fmt.Println("Red Crawler starting...")


	// Extracts the links for each share and the share names
	links, papers := red_crawler.Fundamentus_parser()


	// Extracts the targeted information for every paper
	shares = red_crawler.Crawl_shares(links, papers)


	// Sorts out the 10 most valuable
	top_ten := red_crawler.Top_ten_papers(shares)


	// Creates/revamps the database
	red_crawler.DB_persist(top_ten)
}


//===========================================================================

	// Prints the most valuable shares to the standard output
	/*fmt.Println("\nTop Ten Shares\n\n")
	for i:= 0; i <10; i++ {
		Print_paper(top_ten[i])
	}
	fmt.Println("\n")*/
