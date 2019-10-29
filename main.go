package main

import "encoding/json"
import "os"
import "time"
import "fmt"
import "log"
import "github.com/gocolly/colly"
import "github.com/gocolly/colly/debug"

type Sample struct {
	Hour  string  `json:"Hourly"`
	Value float64 `json:"TotalValue"`
}

const BASE_URL = "https://mywater.redwoodcity.org"

func main() {
	login := colly.NewCollector(colly.Debugger(&debug.LogDebugger{}))
	login.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:70.0) Gecko/20100101 Firefox/70.0"

	csrf := ""
	login.OnHTML("input[name='OuterHeader$hdnCSRFToken']", func(e *colly.HTMLElement) {
		csrf = e.Attr("value")
	})

	login.Visit(BASE_URL + "/Portal/default.aspx")

	login.OnRequest(func(r *colly.Request) {
		r.Headers.Set("csrftoken", csrf)
		r.Headers.Set("Accept", "application/json, text/javascript, */*; q=0.01")
		r.Headers.Set("Content-Type", "application/json;charset=UTF-8")
	})

	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	date := os.Getenv("DATE")
	if len(date) == 0 {
		date = time.Now().AddDate(0, 0, -1).Format("01/02/2006")
	}
	body := []byte(`{"username":"` + username + `","password":"` + password + `","rememberme":true,"calledFrom":"LN"}`)
	err := login.PostRaw(BASE_URL+"/Portal/Default.aspx/validateLogin", body)

	if err != nil {
		log.Fatal(err)
	}

	fetch := login.Clone()
	fetch.OnRequest(func(r *colly.Request) {
		r.Headers.Set("csrftoken", csrf)
		r.Headers.Set("Accept", "application/json, text/javascript, */*; q=0.01")
		r.Headers.Set("Content-Type", "application/json;charset=UTF-8")
	})

	fetch.OnResponse(func(r *colly.Response) {
		var result map[string]string
		json.Unmarshal([]byte(r.Body), &result)

		var table map[string][]Sample
		json.Unmarshal([]byte(result["d"]), &table)

		samples := table["Table"]
		layout := "01/02/2006 3:04 PM"
		for _, sample := range samples {
			t, err := time.Parse(layout, date+" "+sample.Hour)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("water_usage_gal value=%f %d\n", sample.Value, t.UnixNano())
		}
	})

	// Type - 'G' for gallons or 'W' for HCF
	body = []byte(`{"Type":"G","Mode":"H","strDate":"` + date + `","hourlyType":"H","seasonId":0,"weatherOverlay":"0","usageyear":"","MeterNumber":"","BillDate":""}`)
	err = fetch.PostRaw(BASE_URL+"/Portal/Usages.aspx/LoadWaterUsage?type=WU", body)

	if err != nil {
		log.Fatal(err)
	}
}
