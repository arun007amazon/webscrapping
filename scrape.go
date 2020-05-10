package main

import (
	"encoding/csv"
    "log"
    "os"
    "fmt"

	"github.com/gocolly/colly"
)

func main() {
    rows := readSample()
    rows = appendSum(rows)
    fmt.Println(rows)
    writeChanges(rows)
}

func readSample() [][]string {
    f, err := os.Open("data.csv")
    if err != nil {
        log.Fatal(err)
    }
    rows, err := csv.NewReader(f).ReadAll()
    f.Close()
    if err != nil {
        log.Fatal(err)
    }
    return rows
}

func appendSum(rows [][]string) [][]string {

    row := getDatafromUrl(rows[1])
    rows = append(rows, row)
    fmt.Println(rows)
    return rows
}


func getDatafromUrl(urls []string)[]string{
	stockdetails := []string{}
	fmt.Println(urls)
        for _, s := range urls{
		price := "";
		if(s == "time"){
		   stockdetails = append(stockdetails, "time")
		}else{
		   price = getprice(s)
		   stockdetails = append(stockdetails, price)
		}
        }
	return stockdetails
}

func getprice(url string) string{

        // Instantiate the default Collector
        c := colly.NewCollector()

        // Before making a request, print "Visiting ..."
        c.OnRequest(func(r *colly.Request) {
                fmt.Println("Visiting: ", r.URL)
        })

	number := ""
        c.OnHTML(`.nsedata_bx`, func(e *colly.HTMLElement) {
                //Locate and extract different pieces information about each movie
                number = e.ChildText(".span_price_wrap")
                fmt.Println(number)
        })
        // start scraping the page under the given URL
        c.Visit(url)
        fmt.Println("End of scraping: ", url)
 	return number
}

func writeChanges(rows [][]string) {
    f, err := os.Create("data.csv")
    if err != nil {
        log.Fatal(err)
    }
    err = csv.NewWriter(f).WriteAll(rows)
    f.Close()
    if err != nil {
        log.Fatal(err)
    }
}

