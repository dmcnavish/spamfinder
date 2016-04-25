package main

import (
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
)

// Spammer hold info for a single spamming monster 
type Spammer struct {
	From            string
	UnsubscribeInfo string
	UnsubscribeURL  string
	EmailBody       string
}

const UrlRegex = "[-a-zA-Z0-9@:%._\\+~#=]{2,256}\\.[a-z]{2,6}\\b([-a-zA-Z0-9@:%_\\+.~#?&//=]*)" 
    //"(http|ftp|https):\\/\\/([\\w_-]+(?:(?:\\.[\\w_-]+)+))([\\w.,@?^=%&:/~+#-]*[\\w@?^=%&/~+#-])?"

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
		panic(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	r, _ := regexp.Compile(UrlRegex)

	for k, v := range *spammers {
		if err := w.Write([]string{k, parseURL(v.UnsubscribeInfo, r)}); err != nil {
			log.Fatalln("error writing record to csv", err)
		}
	}

	w.Flush()
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
