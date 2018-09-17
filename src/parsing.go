
package main

import (
	"fmt"
	"strings"
	"io/ioutil"
	"net/http"
	"strconv"
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


//===========================================================================


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


// Parses the Fundamentus webpage into lists of links and shares (papers)
func Fundamentus_parser() ([]string, []string) {
	// Takes the HTML of the targeted URL
	html := get_html(TARGET_URL)

	// Parses the HTML for the to the paper names
	papers := trim_select(html, PREFIX_PAPERS, SUFFIX_PAPERS)

	// Builds the string urls to each paper
	links := links_to_papers(papers)

	return links, papers
}



// Extracts the targeted information from an individual URL, prints the result to 'shares[i]'
// Some of the links are broken or no longer available, therefore some error handling is necessary
// Broken links results in such a share declaration: Paper{acronym, "Website error", 0, 0.0}, but execution continues normally
// #REFACTOR: Encapsulate error handling to avoid repetition
func Share_parser(shares []Paper, link string, acronym string, i int, semaphore chan green) {


		// Takes the HTML of a paper
		paper_html := get_html(link)

		// Parses the HTML for the company name
		company := strings.Join(trim_select(paper_html, PREFIX_COMPANY, SUFFIX_COMPANY), ",")
		company, err := strconv.Unquote(`"` + company + `"`)
		if err != nil { panic("Cannot unquote string") }
		if(company == "") {
			shares[i] = Paper{acronym, "Website error", 0, 0.0}
			semaphore <- green{}
			return
		}

		// Parses the HTML for the market value
		helper := strings.Join(strings.Split(trim_select(paper_html, PREFIX_VALUE, SUFFIX_VALUE)[0], "."), "")
		if(helper == "") {
			shares[i] = Paper{acronym, "Website error", 0, 0.0}
			semaphore <- green{}
			return
		}
		value, err := strconv.Atoi(helper)

		// Parses the HTML for the daily fluctuation
		helper = fmt.Sprintf("%s", trim_select(paper_html, PREFIX_FLUCTUATION, SUFFIX_FLUCTUATION))
		if(helper == "") {	
			shares[i] = Paper{acronym, "Website error", 0, 0.0}
			semaphore <- green{}
			return
		}
		helper = strings.Replace(trim_select(helper, PREFIX_FONT, SUFFIX_FONT)[0], ",", ".", 1)
		fluctuation, err := strconv.ParseFloat(helper, 64)
		if err != nil { panic("Cannot convert daily fluctuation to flot64") }

		// Builds a struct for every paper
		shares[i] = Paper{acronym, company, value, fluctuation}

		// Used for DEBUG
		//Print_paper (shares[i])

		// Turn semaphore green
		semaphore <- green{}
}


//===========================================================================
