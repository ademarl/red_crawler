
package main

import (
	"fmt"
	"sort"
)


//===========================================================================


// Main function
func main() {

	fmt.Println("Red Crawler starting...")


	// Extracts the links for each share and the share names
	links, papers := Fundamentus_parser()
	n := len(papers)	


	// Extracts the targeted information for every paper
	shares := Crawl_shares(links, papers, n)


	// Sort out the 10 most valuable
	// #REFACTOR: The number of shares is not that big, therefore, sorting is ok. To optimize: add all elements to a max heap and take the first 10 elements
	sort.SliceStable(shares, func(i, j int) bool { return shares[i].value > shares[j].value })
	top_ten := shares[0 : 10]


	// Creates/revamps the database
	DB_persist(top_ten)
}


//===========================================================================

	// Prints the most valuable shares to the standard output
	/*fmt.Println("\nTop Ten Shares\n\n")
	for i:= 0; i <10; i++ {
		Print_paper(top_ten[i])
	}
	fmt.Println("\n")*/
