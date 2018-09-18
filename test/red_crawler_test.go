
package red_crawler_test

import (
    "testing"
	"red_crawler/src/red_crawler"
)


// #REFACTOR better expession of errors
//===========================================================================


func TestGetHtml(t *testing.T) {

	html := red_crawler.Get_html("https://www.google.com")
	if (html == "") {
		t.Errorf("TestGetHtml failed!")
	}
}

func TestLinks(t *testing.T) {

	input := []string{"ABC", ""}

	link := red_crawler.Links_to_papers(input)
	if (link[0] != red_crawler.PAPER_URL_PREFIX + "ABC") {
		t.Errorf("TestLinks failed!")
	}
}


func TestTrimSelect(t *testing.T) {

	target := red_crawler.Trim_select("header prefixred10suffix", "prefix", "suffix" )
	if (target[0] != "red10") {
		t.Errorf("TestTrimSelect failed!")
	}
}


func TestTopTen(t *testing.T) {

	papers := make([]red_crawler.Paper, 20)

	for i := 0; i < 20; i++ {
		papers[i] = red_crawler.Paper{"", "", i, 10}
	}
	papers = red_crawler.Top_ten_papers(papers)
	for i := 0; i < 10; i++ {
		if (papers[i].Value != 19-i) {
			t.Errorf("TestTopTen failed!")
		}
	}
}


func TestCrawlShares(t *testing.T) {

	papers := []string{"AALR3", "ABCB3"}
	links := []string{"https://www.fundamentus.com.br/detalhes.php?papel=AALR3", "https://www.fundamentus.com.br/detalhes.php?papel=ABCB3"}

	p := red_crawler.Crawl_shares(links, papers)

	red_crawler.Print_paper(p[0])
	red_crawler.Print_paper(p[1])

	if (p[0].Company != "ALLIAR ON NM" || p[1].Company != "ABC Brasil ON N2") {
		t.Errorf("TestCrawlShares failed!")
	}
	if (p[0].Value < 0 || p[1].Value < 0) {
		t.Errorf("TestCrawlShares failed!")
	}
}


//===========================================================================
