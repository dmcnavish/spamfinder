package main

import (
	"bytes"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"golang.org/x/net/html"
	"gopkg.in/xmlpath.v2"
	"log"
	"os"
	"regexp"
	"strings"
)

// Spammer hold info for a single spamming monster
type Spammer struct {
	From            string
	UnsubscribeInfo string
	UnsubscribeURL  string
	EmailBody       string
}

// URLRegex regex to determin if a string is a URL
const URLRegex = "(http|ftp|https):\\/\\/([\\w_-]+(?:(?:\\.[\\w_-]+)+))([\\w.,@?^=%&:/~+#-]*[\\w@?^=%&/~+#-])?"

func getSpammers() *map[string]*Spammer {
	srv := GetService()
	query := "is:inbox unsubscribe after:2016/01/01"
	msgs, err := Search(srv, query)
	if err != nil || msgs == nil {
		log.Fatal("failed to get Messages", err)
	}
	fmt.Printf("total msgs: %v\n", len(msgs))
	spammers := make(map[string]*Spammer)
	for _, m := range msgs {
		s := &Spammer{}
		for _, h := range m.Payload.Headers {
			if h.Name == "List-Unsubscribe" {
				s.UnsubscribeInfo = h.Value
			}
			if h.Name == "From" {
				s.From = h.Value
			}
		}
		ds, err := base64.URLEncoding.DecodeString(m.Payload.Body.Data)
		if err != nil {
			log.Fatal("failed to decode message", err)
		}
		s.EmailBody = string(ds)
		spammers[s.From] = s
	}
	return &spammers
}

func exportSpammers(spammers *map[string]*Spammer, outputFilename string) {
	f, err := os.Create(outputFilename)
	if err != nil {
		fmt.Println("Error creating output file. File is probably open")
		panic(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)

	for k, v := range *spammers {
		url := findUnsubscribeURL(v)
		if err := w.Write([]string{k, url}); err != nil {
			log.Fatalln("error writing record to csv", err)
		}
	}

	w.Flush()
}

func findUnsubscribeURL(spammer *Spammer) string {
	r, _ := regexp.Compile(URLRegex)
	url := parseURL(spammer.UnsubscribeInfo, r)
	if url == "" {
		url = parseWithXpath(spammer.EmailBody)
	}

	return url
}

func parseWithXpath(text string) string {
	path := xmlpath.MustCompile("//a[@title='Unsubscribe' or @title='unsubscribe' or text() = 'unsubscribe' or text() = 'Unsubscribe']/@href")
	text = cleanString(text)
	r := strings.NewReader(text)
	root, err := xmlpath.Parse(r)
	if err != nil {
		fmt.Printf("error using xpath for text: %v \n %v \n", text, err)
		return ""
	}
	if value, ok := path.String(root); ok {
		fmt.Println("found: ", value)
		return value
	}
	return ""
}

func cleanString(text string) string {
	//We need to rerender the text this way because gopkg.in/xmlpath.v2 fails if there are any errors in the HTML
	reader := strings.NewReader(text)
	root, err := html.Parse(reader)
	if err != nil {
		fmt.Printf("Error parsing text: %v \n %v \n", text, err)
		return text
	}

	var b bytes.Buffer
	html.Render(&b, root)
	return b.String()
}

func parseURL(s string, r *regexp.Regexp) string {
	match := r.FindString(s)
	if match == "" {
		fmt.Printf("Found no match: %s\n", s)
	}
	return match
}

func main() {
	spammers := getSpammers()
	outputFilename := "output.csv"
	exportSpammers(spammers, outputFilename)
}
