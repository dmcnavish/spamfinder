package main 

import (
    "testing"
    "regexp"
) 

var validUrls = []string{
    "http://test.com?unsubscribe=true",
    "https://www.golang-book.com/books/intro/12",
    "http://test.com?unsubscribe=true&something=false",
    "www.test.com",
    "www.test.com?unsubscribe=true",
}

var invalidUrls = []string{
    "not_a_url_com",
    "htp://no",
    "??something??",
}

func TestParseURL_match(t *testing.T){
    r, _ := regexp.Compile(UrlRegex)
    for _ , url := range validUrls {
        result := parseURL(url, r)
        if result == "" {
            t.Error("URL not found!", result)
        }        
    }
}

func TestParseURL_noMatch(t *testing.T){
    r, _ := regexp.Compile(UrlRegex)
    for _, url := range invalidUrls {
        result := parseURL(url, r)
        if result != "" {
            t.Error("No match should have been found", result)
        }    
    }
}

func testExportSpammers(t *testing.T){
    
}