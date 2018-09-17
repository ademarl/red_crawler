
package main

import "fmt"


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
