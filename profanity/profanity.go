// Package profanity provides basic functions for checking profanity in a text.
// Its primary job is to filter profanity from text with high level of concurrency.
package profanity

import (
	"io/ioutil"
	"strings"
	"sync"
	"regexp"
	"os"
)

var (
	// BlankRegexp contains pattern that will be replaced with blank character.
	BlankRegexp = regexp.MustCompile("([\\[$&,:;=?#|'<>.^*\\(\\)%\\]])|(\\b\\d+\\b)|(cum\u0020laude)|(he\\'ll)|(\\B\\#)|(&\\#?[a-z0-9]{2,8};)|(\\b\\'+)|(\\'+\\b)|(\\b\\\")|(\\\"\\b)|(dick\u0020cheney)|(\\!+\\B)")
	// URepeatRegexp contains regex for continuous occurrence of u.
	URepeatRegexp = regexp.MustCompile("u+")
	// IRepeatRegexp contains regex for continuous occurrence of i.
	IRepeatRegexp = regexp.MustCompile("i+")
)

// wordsMap is the map containing profane words as map key.
var wordsMap map[string]interface{} = make(map[string]interface{})

// init cache the profane words.
func init() {
	cacheAbuses()
}

// Profanity contains the result of any profanity check.
type Profanity struct {
	Total int
	Found []string
}

// Find finds the profanity in a text.
// Find returns Profanity.
func Find(txt string) Profanity {
	channel := make(chan string)
	var wg sync.WaitGroup
	found := []string{}
	filterdText := filterUsingRegex(txt)
	words := strings.Split(filterdText, " ")
	wg.Add(len(words))
	go func() {
		for msg := range channel {
			found = append(found, msg)
			wg.Done()
		}
	}()
	for _, word := range words {
		go func(w string) {
			s := strings.TrimSpace(w);
			if _, ok := wordsMap[s]; s != "" && ok {
				channel <- s
			} else {
				wg.Done()
			}
		}(word)
	}
	wg.Wait()
	close(channel)
	return Profanity{Total:len(found), Found:found}
}

// Check checs for profanity in a text.
// Check returns true or false based on profanity in text.
func Check(txt string) bool {
	return Find(txt).Total > 0
}

// filterUsingRegex filters the text using regex and replaces text with appropriate character.
// filterUsingRegex uses BlankRegexp, URepeatRegexp and IRepeatRegexp to find a text match.
func filterUsingRegex(text string) string {
	text = BlankRegexp.ReplaceAllString(text, "")
	text = URepeatRegexp.ReplaceAllString(text, "u")
	text = IRepeatRegexp.ReplaceAllString(text, "i")
	return text
}

// cacheAbuses cache the data of directory profanity/data into wordsMap.
func cacheAbuses() {
	CacheDirContent("profanity/data")
}

// CacheDirContent cache the data of all files in specified directory into wordsMap.
func CacheDirContent(dir string) {
	_, err := os.Stat(dir)
	if err == nil {
		files, err := ioutil.ReadDir(dir)
		checkErr(err)
		if len(files) > 0 {
			for _, file := range files {
				if file.Mode().IsRegular() && !file.IsDir() {
					f, err := ioutil.ReadFile(dir + "/" + file.Name())
					checkErr(err)
					if err == nil {
						words := strings.Split(string(f), "\n")
						for _, s := range words {
							wordsMap[strings.TrimSpace(s)] = nil
						}
					}
				}
			}
		}
	}
}

// checkErr checks for err.
func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

