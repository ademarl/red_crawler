package main

import (
	"fmt"
	_"math/big"
	"sort"
	"strings"
	"io/ioutil"
	"net/http"
	"strconv"
	_ "database/sql"
    _ "github.com/go-sql-driver/mysql"
)


// Constant strings used for the URL and splicing the HTML files
const targetURL = "https://www.fundamentus.com.br/detalhes.php"
const paperUrlPrefix = "https://www.fundamentus.com.br/detalhes.php?papel="
const prefix_papers = "detalhes.php?papel="
const suffix_papers = "\""
const prefix_company = "Empresa</span></td>\n					<td class=\"data\"><span class=\"txt\">"
const suffix_company = "<"
const prefix_value = "Valor de mercado</span></td>\n					<td class=\"data w3\"><span class=\"txt\">"
const suffix_value = "<"
const prefix_fluctuation = "Dia</span></td>\n					<td class=\"data w1\"><span class=\"oscil\"><font color="
const suffix_fluctuation = "<"
const prefix_font = ">"
const suffix_font = "%"

type Paper struct {

	share string			// acronym
	company string			// company name
	value int			// market value
	fluctuation float64		// daily fluctuation
}


func main() {
	// fmt.Println("Red Crawler starting...")











	// Takes the HTML of the targeted URL
	html := get_html(targetURL)

	// Trims the HTML for the to the paper names
	papers := trim_select(html, prefix_papers, suffix_papers)
	n := len(papers)

	// Builds the string urls to each paper
	links := links_to_papers(papers)

	// Extracts the targeted information for every paper
	shares := crawl_shares(links, papers, n)

	fmt.Println("\n\n\n")
	for i:= 0; i <12; i++ {
		printPaper(shares[i])
	}

	// Sort out the 10 most valuable
	// REFACTOR: The number of shares is not that big, therefore, sorting is ok. To optimize: add all elements to a max heap and take the first 10 elements
	sort.SliceStable(shares, func(i, j int) bool { return shares[i].value > shares[j].value })
	top_ten := shares[0 : 10]

/*
	fmt.Println("\nTop Ten Shares\n\n")
	for i:= 0; i <10; i++ {
		printPaper(top_ten[i])
	}
	fmt.Println("\n")*/

	// Create a database and insert the 10 itens










}


// Auxiliar function to print the targeted elements of a paper
func printPaper (share Paper) {
	fmt.Println(share.share, share.company, share.value, share.fluctuation)
}



// Returns the html of a website
func get_html(url string) string {

	// HTTP request to the targeted website
	resp, err := http.Get(url)
	if err != nil { panic("Cannot read website") }

	// Extract the HTML
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil { panic("Cannot real html") }

	// Returns as a string
	return fmt.Sprintf("%s", html)
}

// Trim the HTML to a targeted section of text
func trim_select(targetHTML string, prefix string, suffix string) []string {

	// Makes a slice by cutting the prefixes
	var papers = strings.Split(targetHTML, prefix)

	// Makes another slice by cutting the suffixes
	for i := 1; i < len(papers); i++ {
		papers[i] = strings.Split(papers[i], suffix)[0]
	}

	// Ignore the first slice that appears before any targets
	return papers[1 : len(papers)]
}


//
func links_to_papers(papers []string) []string {

	var links = make([]string, len(papers))

	for i := 0; i < len(papers); i++ {
		links[i] = paperUrlPrefix + papers[i]
	}

	return links
}

func crawl_shares(links []string, papers []string, n int) []Paper {

	var shares = make([]Paper, n)
	var paper_html = make([]string, n)
	for i := 0; i < n; i++ {


		// Takes the HTML of a paper
		paper_html[i] = get_html(links[i])

		// Trims the HTML for the company name, market value and daily fluctuation
		company := strings.Join(trim_select(paper_html[i], prefix_company, suffix_company), ",")
		if(company == "") {
			shares[i] = Paper{papers[i], "Website error", 0, 0.0}
			//printPaper (shares[i])
			continue
		}

		helper := strings.Join(strings.Split(trim_select(paper_html[i], prefix_value, suffix_value)[0], "."), "")
		if(helper == "") {
			shares[i] = Paper{papers[i], "Website error", 0, 0.0}
			//printPaper (shares[i])
			continue
		}
		value, err := strconv.Atoi(helper)

		helper = fmt.Sprintf("%s", trim_select(paper_html[i], prefix_fluctuation, suffix_fluctuation))
		if(helper == "") {
			shares[i] = Paper{papers[i], "Website error", 0, 0.0}
			//printPaper (shares[i])
			continue
		}
		helper = strings.Replace(trim_select(helper, prefix_font, suffix_font)[0], ",", ".", 1)
		fluctuation, err := strconv.ParseFloat(helper, 64)
		if err != nil { panic("Cannot convert daily fluctuation to flot64") }

		// Builds a struct for every paper
		shares[i] = Paper{papers[i], company, value, fluctuation}

		//printPaper (shares[i])
	}
	return shares
}




































