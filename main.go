package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/klauspost/compress/gzhttp"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	l              = log.New(os.Stdout, "[CIQ-ZendeskCrawler] ", 0)
	responseStruct = responseBody{}
)

func main() {
	if len(os.Args) != 4 {
		l.Printf("Username and password args must be specified as such: [ ./crawler 'johndoe@ciq.com' 'myPassword' /path/to/urls.txt ]")
		os.Exit(1)
	}
	username := os.Args[1]
	password := os.Args[2]
	pathList := os.Args[3]

	getBodies(username, password, pathList)
}

func getBodies(username string, password string, pathList string) {
	// Setup an HTTP client with a large buffer & GZIP transport
	transport := &http.Transport{
		WriteBufferSize: 5000000,
		ReadBufferSize:  5000000,
	}
	client := &http.Client{
		Transport: gzhttp.Transport(transport),
	}

	f, err := os.Open(pathList)
	defer f.Close()
	if err != nil {
		l.Printf("Error opening url list file: %s\n", err.Error())
		return
	}

	urlScanner := bufio.NewScanner(f)

	//buf := bytes.NewBuffer(make([]byte, 1000000))

	var req *http.Request
	var resp *http.Response
	var decoder *json.Decoder

	var f2 *os.File
	defer f2.Close()

	for urlScanner.Scan() {
		if urlScanner.Text() == "" || urlScanner.Text() == "\n" {
			break
		}
		l.Printf("Sending request to %s\n", urlScanner.Text())
		req, err = http.NewRequest("GET", urlScanner.Text(), nil)
		if err != nil {
			l.Printf("Error generating new GET request: %s\n", err.Error())
			continue
		}
		req.Header.Set("Authorization", "Basic"+basicAuth(username, password))
		resp, err = client.Do(req)
		if err != nil {
			l.Printf("Error sending new GET request: %s\n", err.Error())
			continue
		}
		decoder = json.NewDecoder(resp.Body)
		err = decoder.Decode(&responseStruct)
		if err != nil {
			l.Printf("Error decoding response body: %s\n", err.Error())
			continue
		}

		for i := range responseStruct.Articles {
			rp := responseStruct.Articles[i]
			err = os.Mkdir(strings.ReplaceAll(rp.Title, " ", ""), 0766)
			if err != nil {
				err = os.Mkdir(strings.ReplaceAll(rp.Title+"-2", " ", ""), 0766)
				if err != nil {
					l.Printf("Error writing fallback file: %s\n", err.Error())
					continue
				}
			}
			err = ioutil.WriteFile(strings.ReplaceAll(rp.Title, " ", "")+"/index.html", []byte(rp.Body), 0766)
			if err != nil {
				l.Printf("Error writing HTML body to file: %s\n", err.Error())
				continue
			}
			// Adds a comma separator for JSON arrays, there is probably a better way to do this(?)
			labelNames := "'" + strings.Join(rp.LabelNames, `", "`) + `'`

			metadata := fmt.Sprintf(`
{
	"Metadata": {
		"AuthorID": %d,
		"CommentsDisabled": %t,
		"CreatedAt": "%s",
		"Draft": %t,
		"EditedAt": "%s",
		"HtmlURL": "%s",
		"ID": %d,
		"LabelNames": "[%s]",
		"Locale": "%s",
		"Name": "%s",
		"Outdated": %t,
		"PermissionGroupID": %d,
		"Position": %d,
		"Promoted": %t,
		"SectionID": %d,
		"SourceLocale": "%s",
		"Title": "%s",
		"UpdatedAt": "%s",
		"URL": "%s",
		"UserSegmentID": %d,
		"VoteCount": %d,
		"VoteSum": %d
	}
}
`, rp.AuthorID, rp.CommentsDisabled, rp.CreatedAt, rp.Draft, rp.EditedAt,
				rp.HTMLURL, rp.ID, labelNames, rp.Locale, rp.Name, rp.Outdated, rp.PermissionGroupID,
				rp.Position, rp.Promoted, rp.SectionID, rp.SourceLocale, rp.Title, rp.UpdatedAt, rp.URL,
				rp.UserSegmentID, rp.VoteCount, rp.VoteSum)

			err = ioutil.WriteFile(strings.ReplaceAll(rp.Title, " ", "")+"/metadata.json", []byte(metadata), 0766)
			l.Printf("Wrote %s\n", strings.ReplaceAll(rp.Title, " ", "")+"/index.html")
		}

	}

	err = resp.Body.Close()
	if err != nil {
		l.Printf("Error closing response body: %s\n", err.Error())
	}

}
func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
