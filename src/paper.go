
package main

import "fmt"


//===========================================================================
// Defines a struct for a stock market share with the 4 relevant fields
// Also defines a print method for the struct
//===========================================================================

// A struct representing a share with the relevant data
type Paper struct {

	share string			// acronym
	company string			// company name
	value int				// market value
	fluctuation float64		// daily fluctuation
}


// Auxiliar function to print the targeted elements of a share
func Print_paper (share Paper) {

	fmt.Println(share.share, share.company, share.value, share.fluctuation)
}


//===========================================================================
