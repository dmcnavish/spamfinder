package main

import (
	"regexp"
	"testing"
)

var validUrls = []string{
	"http://test.com?unsubscribe=true",
	"https://www.golang-book.com/books/intro/12",
	"http://test.com?unsubscribe=true&something=false",
	"http://www.test.com",
	"http://www.test.com?unsubscribe=true",
}

var invalidUrls = []string{
	"not_a_url_com",
	"htp://no",
	"??something??",
}

var validXpathStrings = []string{
	`<People>
        <Person>
            <FullName>Jerome Anthony</FullName>
        </Person>
        <Person>
            <a href="http://x.e.flyfrontier.com/ats/show.aspx?cr=551&amp;fm=19&amp;tp=i-H55-8t-Zk-11ideA-1r-Xbhr-1c-BWM-11W2GI-9ol21" title="Unsubscribe" target="_blank"><font color="#00acec">Unsubscribe</font></a>
            <FullName>Christina</FullName>
        </Person>
     </People>`,
	`<html>
        <body>
            <p>
                <a href="http://unsubscripe.com" title="unsubscribe" target="_blank"><font>Unsubscribe</font></a>
            </p>
        </body>
    </html>`,
	`<p>
        &nbsp;<a href="http://unsubscripe.com" title="unsubscribe" target="_blank"><font>Unsubscribe</font></a>
    </p>
    `,
	`<p>
        <a href="http://unsubscripe.com">Unsubscribe</a>
    </p>
    `,
	`<td width="100%" bgcolor="#F1F1F1" style="padding-top:0;padding-right:20px;padding-bottom:25px;padding-left:20px;font-family:Arial,Helvetica,sans-serif;font-size:11px;line-height:13pt;color:#333333;text-align:left">
        <div style="min-height:20px">&nbsp;</div>
        <a href="http://info.pivotal.io/n0Ii0LK02qAN0JU08i00Cs3" style="text-decoration:none;color:#6db33f" target="_blank">Unsubscribe</a>.</td>
    `,
}

var invalidXpathStrings = []string{
	`<html>
        <body>
            <p>
                <a href="http://unsubscripe.com" target="_blank"><font>Unsubscribe</font></a>
            </p>
         </body>
     </html>`,
	`<html>
        <body>
            <p></p>
        </body>
    </html>`,
}

func TestParseURL_match(t *testing.T) {
	r, _ := regexp.Compile(URLRegex)
	for _, url := range validUrls {
		result := parseURL(url, r)
		if result == "" {
			t.Error("URL not found!", result)
		}
	}
}

func TestParseURL_noMatch(t *testing.T) {
	r, _ := regexp.Compile(URLRegex)
	for _, url := range invalidUrls {
		result := parseURL(url, r)
		if result != "" {
			t.Error("No match should have been found", result)
		}
	}
}

func TestParseWithXpath(t *testing.T) {
	for _, text := range validXpathStrings {
		result := parseWithXpath(text)
		r, _ := regexp.Compile(URLRegex)
		if parseURL(result, r) == "" {
			t.Error("No URL found for text: ", text)
		}
	}
}

func TestParseWithXpath_noMatch(t *testing.T) {
	for _, text := range invalidXpathStrings {
		result := parseWithXpath(text)
		r, _ := regexp.Compile(URLRegex)
		if parseURL(result, r) != "" {
			t.Error("Expected no match for text: ", text)
		}
	}
}
