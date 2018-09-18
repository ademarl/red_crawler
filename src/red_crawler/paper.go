
package red_crawler

import(
		"fmt"
		"sort"
)

//===========================================================================
// Defines a struct for a stock market share with the 4 relevant fields
// Also defines a print method for the struct
//===========================================================================

// A struct representing a share with the relevant data
type Paper struct {

	Share string			// acronym
	Company string			// company name
	Value int				// market value
	Fluctuation float64		// daily fluctuation
}


// Sorts and returns the 10 higher market value Papers
// #REFACTOR: The number of shares is not that big, therefore, sorting is ok. To optimize: add all elements to a max heap and take the first 10 elements
func Top_ten_papers (shares []Paper) []Paper{

	sort.SliceStable(shares, func(i, j int) bool { return shares[i].Value > shares[j].Value })
	return shares[0 : 10]
}

// Auxiliar function to print the targeted elements of a share
func Print_paper (share Paper) {

	fmt.Println(share.Share, share.Company, share.Value, share.Fluctuation)
}


//===========================================================================
