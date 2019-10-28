package main

import "os"
import "fmt"
import "github.com/gocolly/colly"
import "github.com/gocolly/colly/debug"

func main() {
	login := colly.NewCollector(colly.Debugger(&debug.LogDebugger{}))
	login.UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:70.0) Gecko/20100101 Firefox/70.0"

	//  curl -X POST https://mywater.redwoodcity.org/Portal/Usages.aspx/LoadWaterUsage\?type=WU --verbose -H "Cookie: ASP.NET_SessionId=aydvbc3krdza0tophcfqbkyx;" -H "Accept: application/json, text/javascript, */*; q=0.01" -H "csrftoken: hRxoQG/KSDC6/WHCpWC0Pydi75PJujcRsZdlXtrLiJaa94n3js+TgTVmaiGGKZes" -d '{"Type":"W","Mode":"H","strDate":"10/26/2019","hourlyType":"H","seasonId":0,"weatherOverlay":"0","usageyear":"","MeterNumber":"","BillDate":""}' -H "Content-Type: application/json"

	csrf := ""
	login.OnHTML("input[name='OuterHeader$hdnCSRFToken']", func(e *colly.HTMLElement) {
		csrf = e.Attr("value")
		fmt.Printf("%s\n", csrf)
	})

	login.Visit("https://mywater.redwoodcity.org/Portal/default.aspx")

	login.OnResponse(func(r *colly.Response) {
		siteCookies := login.Cookies(r.Request.URL.String())
		fmt.Printf("%+v\n", siteCookies)
	})

	login.OnRequest(func(r *colly.Request) {
		r.Headers.Set("csrftoken", csrf)
		r.Headers.Set("Accept", "application/json, text/javascript, */*; q=0.01")
		r.Headers.Set("Content-Type", "application/json;charset=UTF-8")
	})

	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	err := login.PostRaw("https://mywater.redwoodcity.org/Portal/Default.aspx/validateLogin", []byte(`{"username":"`+username+`","password":"`+password+`","rememberme":true,"calledFrom":"LN"}`))

	if err != nil {
		fmt.Println(err)
	}

	fetch := login.Clone()
	fetch.OnRequest(func(r *colly.Request) {
		r.Headers.Set("csrftoken", csrf)
		r.Headers.Set("Accept", "application/json, text/javascript, */*; q=0.01")
		r.Headers.Set("Content-Type", "application/json;charset=UTF-8")
	})

	fetch.OnResponse(func(r *colly.Response) {
		fmt.Printf("%s\n", r.Body)
	})

	err = fetch.PostRaw("https://mywater.redwoodcity.org/Portal/Usages.aspx/LoadWaterUsage?type=WU", []byte(`{"Type":"W","Mode":"H","strDate":"10/26/2019","hourlyType":"H","seasonId":0,"weatherOverlay":"0","usageyear":"","MeterNumber":"","BillDate":""}`))

	if err != nil {
		fmt.Println(err)
	}
}
