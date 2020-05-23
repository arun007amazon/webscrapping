package main

import (
    "encoding/csv"
    "log"
    "os"
    "fmt"
    "github.com/gocolly/colly"
    "bytes"
    "encoding/json"
    "errors"
    "net/http"
    "time"
    "strconv"
)

type SlackRequestBody struct {
    Text string `json:"text"`
}

type Configuration struct {
    Webhook   string
}


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
    sendDataToSlack(rows[0], row)
    fmt.Println(rows)
    return rows
}


func getDatafromUrl(urls []string)[]string{
	stockdetails := []string{}
	fmt.Println(urls)
	i := time.Now().Unix()
	unixTime := strconv.FormatInt(i, 10)
        for _, s := range urls{
		price := "";
		if(s == "time"){
		   stockdetails = append(stockdetails, unixTime)
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

func sendDataToSlack(companys []string, todaysPrice []string){

    //get config file information
    file, _ := os.Open("slackconfig.json")
    defer file.Close()
    decoder := json.NewDecoder(file)
    configuration := Configuration{}
    err := decoder.Decode(&configuration)
    if err != nil {
      fmt.Println("error:", err)
    }
    fmt.Println(configuration.Webhook)


    data :=""
    webhookUrl := configuration.Webhook

    //Get local time and append it to slack timestamp
    loc, _ := time.LoadLocation("Asia/Kolkata")
    t :=time.Now().In(loc)
    data = data + "*" + companys[0] + "*: " +t.String() +"\n"

    for i := 1; i < len(companys); i++ {
	data = data + "*" + companys[i] + "*: " +todaysPrice[i] +"\n"
    }
    fmt.Println("data :", data)
    err = SendSlackNotification(webhookUrl, data)
    if err != nil {
        log.Fatal(err)
    }
}

func SendSlackNotification(webhookUrl string, msg string) error {

    slackBody, _ := json.Marshal(SlackRequestBody{Text: msg})
    req, err := http.NewRequest(http.MethodPost, webhookUrl, bytes.NewBuffer(slackBody))
    if err != nil {
        return err
    }

    req.Header.Add("Content-Type", "application/json")

    client := &http.Client{Timeout: 10 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }

    buf := new(bytes.Buffer)
    buf.ReadFrom(resp.Body)
    if buf.String() != "ok" {
        return errors.New("Non-ok response returned from Slack")
    }
    return nil
}
