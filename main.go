package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly/v2"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Item struct {
	Author  string   `json:"author"`
	Time    string   `json:"time"`
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Img     []string `json:"img"`
	Oo      int      `json:"oo"`
	Xx      int      `json:"xx"`
}
type PageInfo struct {
	StatusCode int
	Items      []Item
}

var SCache = &SimpleCache{}

func handler(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	date := r.URL.Query().Get("date")
	page := r.URL.Query().Get("page")
	if page == "" {
		page = "158"
	}
	if date == "" {
		date = now.Format("20060102")
	}

	param := base64.StdEncoding.EncodeToString([]byte(date + "-" + page))

	c := colly.NewCollector()

	p := &PageInfo{Items: make([]Item, 0)}

	// count links
	c.OnHTML("li", func(e *colly.HTMLElement) {
		link := e.Attr("id")
		author := e.DOM.Find("b").Text()
		time := e.DOM.Find("span[class=\"time\"]").Text()
		time = strings.ReplaceAll(time, "@ ", "")
		//otext:=e.DOM.Find("span[class=\"tucao-unlike-container\"] span").Text()
		commentText := e.DOM.Find("p").Text()
		img := make([]string, 0)
		e.ForEach("a[class=\"view_img_link\"]", func(i int, element *colly.HTMLElement) {
			img = append(img, element.Attr("href"))
		})
		oo := 0
		xx := 0
		e.ForEach("span[class=\"tucao-unlike-container\"] span", func(i int, element *colly.HTMLElement) {
			if i == 0 {
				oo, _ = strconv.Atoi(element.Text)

			} else {
				xx, _ = strconv.Atoi(element.Text)
			}
		})
		p.Items = append(p.Items, Item{Title: link, Author: author, Oo: oo, Xx: xx, Time: time, Content: commentText, Img: img})
	})

	// extract status code
	c.OnResponse(func(r *colly.Response) {
		log.Println("response received", r.StatusCode)
		p.StatusCode = r.StatusCode
	})
	c.OnError(func(r *colly.Response, err error) {
		log.Println("error:", r.StatusCode, err)
		p.StatusCode = r.StatusCode
	})

	c.Visit(fmt.Sprintf("http://i.jandan.net/pic/%s#comments", param))

	// dump results
	b, err := json.Marshal(p)
	if err != nil {
		log.Println("failed to serialize response:", err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Write(b)
}

func main() {
	addr := ":7171"
	http.HandleFunc("/", handler)
	log.Println("listening on", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
