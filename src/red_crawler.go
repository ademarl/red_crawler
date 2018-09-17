package main

import (
	"fmt"
	_"math/big"
	"sort"
	"strings"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	_ "database/sql"
    _ "github.com/go-sql-driver/mysql"
)


//===========================================================================


// Constant strings used for the URL and splicing the HTML files
// Suffixes and prefixes were determined by manually searching the HTML files for patterns
const TARGET_URL = "https://www.fundamentus.com.br/detalhes.php"
const PAPER_URL_PREFIX = "https://www.fundamentus.com.br/detalhes.php?papel="
const PREFIX_PAPERS = "detalhes.php?papel="
const SUFFIX_PAPERS = "\""
const PREFIX_COMPANY = "Empresa</span></td>\n					<td class=\"data\"><span class=\"txt\">"
const SUFFIX_COMPANY = "<"
const PREFIX_VALUE = "Valor de mercado</span></td>\n					<td class=\"data w3\"><span class=\"txt\">"
const SUFFIX_VALUE = "<"
const PREFIX_FLUCTUATION = "Dia</span></td>\n					<td class=\"data w1\"><span class=\"oscil\"><font color="
const SUFFIX_FLUCTUATION = "<"
const PREFIX_FONT = ">"
const SUFFIX_FONT = "%"

// Wait time between requests to the target website in milliseconds. Might need calibration
const WAIT = 50

// A struct representing a share with the relevant data
type Paper struct {

	share string			// acronym
	company string			// company name
	value int				// market value
	fluctuation float64		// daily fluctuation
}

// Green light for semaphores used in the crawl_shares function
type green struct {}


//===========================================================================


// Main function
func main() {
	fmt.Println("Red Crawler starting...")

	// Takes the HTML of the targeted URL
	html := get_html(TARGET_URL)

	// Parses the HTML for the to the paper names
	papers := trim_select(html, PREFIX_PAPERS, SUFFIX_PAPERS)
	n := len(papers)

	// Builds the string urls to each paper
	links := links_to_papers(papers)

	// Extracts the targeted information for every paper
	shares := crawl_shares(links, papers, n)

	// Sort out the 10 most valuable
	// REFACTOR: The number of shares is not that big, therefore, sorting is ok. To optimize: add all elements to a max heap and take the first 10 elements
	sort.SliceStable(shares, func(i, j int) bool { return shares[i].value > shares[j].value })
	top_ten := shares[0 : 10]

	// Prints the most aluable shares to the standard output
	fmt.Println("\nTop Ten Shares\n\n")
	for i:= 0; i <10; i++ {
		printPaper(top_ten[i])
	}
	fmt.Println("\n")

	// Create a database and insert the 10 itens
}


//===========================================================================


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


// Takes the name of each share and appends the PAPER_URL_PREFIX to create a full URL for each share
func links_to_papers(papers []string) []string {

	var links = make([]string, len(papers))

	for i := 0; i < len(papers); i++ {
		links[i] = PAPER_URL_PREFIX + papers[i]
	}

	return links
}


// Extracts the targeted information from all URLs by crawling through the links with parallel requests
func crawl_shares(links []string, papers []string, n int) []Paper {

	var shares = make([]Paper, n)

	// Parallel crawling with channels as semaphores
	semaphore := make(chan green, n);
	for i := range papers {
		go share_parser(shares, links[i], papers[i], i, semaphore)
		// Server is unhappy with too many requests, slow it down
		time.Sleep(WAIT*time.Millisecond)
	}

	// Wait until all semaphores are green
	for i := 0; i < n; i++ { <-semaphore }

	return shares
}


//===========================================================================


// Extracts the targeted information from an individual URL, prints the result to 'shares[i]'
// Some of the links are broken or no longer available, therefore some error handling is necessary
// Broken links results in such a share declaration: Paper{acronym, "Website error", 0, 0.0}, but execution continues normally
// REFACTOR: Encapsulate error handling to avoid repetition
func share_parser(shares []Paper, link string, acronym string, i int, semaphore chan green) {


		// Takes the HTML of a paper
		paper_html := get_html(link)

		// Parses the HTML for the company name
		company := strings.Join(trim_select(paper_html, PREFIX_COMPANY, SUFFIX_COMPANY), ",")
		if(company == "") {
			shares[i] = Paper{acronym, "Website error", 0, 0.0}
			//printPaper (shares[i])
			semaphore <- green{}
			return
		}

		// Parses the HTML for the market value
		helper := strings.Join(strings.Split(trim_select(paper_html, PREFIX_VALUE, SUFFIX_VALUE)[0], "."), "")
		if(helper == "") {
			shares[i] = Paper{acronym, "Website error", 0, 0.0}
			//printPaper (shares[i])
			semaphore <- green{}
			return
		}
		value, err := strconv.Atoi(helper)

		// Parses the HTML for the daily fluctuation
		helper = fmt.Sprintf("%s", trim_select(paper_html, PREFIX_FLUCTUATION, SUFFIX_FLUCTUATION))
		if(helper == "") {	
			shares[i] = Paper{acronym, "Website error", 0, 0.0}
			//printPaper (shares[i])
			semaphore <- green{}
			return
		}
		helper = strings.Replace(trim_select(helper, PREFIX_FONT, SUFFIX_FONT)[0], ",", ".", 1)
		fluctuation, err := strconv.ParseFloat(helper, 64)
		if err != nil { panic("Cannot convert daily fluctuation to flot64") }

		// Builds a struct for every paper
		shares[i] = Paper{acronym, company, value, fluctuation}

		// Used for DEBUG
		//printPaper (shares[i])

		// Turn semaphore green
		semaphore <- green{}
}






























