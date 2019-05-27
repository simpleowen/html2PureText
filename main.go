package main

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"strings"
	"io/ioutil"
)

func getBody(doc *html.Node) string {
	var pureText string
	var extractTextNode func(*html.Node)

	extractTextNode = func(n *html.Node) {
		if n.Type == html.TextNode && n.Parent.Data == "span" {
			// fmt.Printf("%s", strings.TrimSpace(strings.Replace(n.Data, "\n", " ", -1)))
			// pureText += strings.TrimSpace(strings.Replace(n.Data, "\n", " ", -1))
			r := strings.NewReplacer("\u00a0","", "\n","", "？", "?", "。",".", "：",":", "，",",","；",";", "（","(", "）",")", "、",",", " ","", "／","/", "“","\"","”","\"", "×","X", "！","!")
			pureText += strings.TrimSpace(r.Replace(n.Data))
			// pureText += strings.TrimSpace(strings.Replace(strings.Replace(n.Data, "\u00a0", "", -1), "\n", " ", -1))
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			extractTextNode(c)
		}
	}

	extractTextNode(doc)

	return pureText

}

func main() {

	parseHTML := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			fmt.Println("准备解析参数")

			htmlText, err := ioutil.ReadAll(r.Body)

			if err != nil {
				fmt.Fprintf(w, "Error occured.")
			}
			// fmt.Println(string(htmlText))

			htmlNode, err := html.Parse(strings.NewReader(string(htmlText)))
			if err != nil {
				fmt.Fprintf(w, "Error occured.")
			}
			pureTEXT := getBody(htmlNode)
			fmt.Println(pureTEXT)
			fmt.Fprintf(w, pureTEXT)

		} else {
			fmt.Fprintf(w, "Please use POST method")
		}
	}

	http.HandleFunc("/html", parseHTML)
	http.ListenAndServe(":8080", nil)
}
