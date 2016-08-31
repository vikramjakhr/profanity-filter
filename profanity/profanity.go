package profanity

import (
	"io/ioutil"
	"strings"
	"sync"
	"regexp"
	"os"
)

const (
	BlankRegexp = regexp.MustCompile("([\\[$&,:;=?#|'<>.^*\\(\\)%\\]])|(\\b\\d+\\b)|(cum\u0020laude)|(he\\'ll)|(\\B\\#)|(&\\#?[a-z0-9]{2,8};)|(\\b\\'+)|(\\'+\\b)|(\\b\\\")|(\\\"\\b)|(dick\u0020cheney)|(\\!+\\B)")
	URepeatRegexp = regexp.MustCompile("u+")
	IRepeatRegexp = regexp.MustCompile("i+")
)

var wordsMap map[string]interface{} = make(map[string]interface{})

func init() {
	cacheAbuses()
}

type Profanity struct {
	Total int
	Found []string
}

func Find(txt string) (Profanity) {
	channel := make(chan string)
	var wg sync.WaitGroup
	found := []string{}
	b, err := ioutil.ReadAll(txt)
	filterdText := filterUsingRegex(string(b))
	checkErr(err)
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

func filterUsingRegex(text string) string {
	text = BlankRegexp.ReplaceAllString(text, "")
	text = URepeatRegexp.ReplaceAllString(text, "u")
	text = IRepeatRegexp.ReplaceAllString(text, "i")
	return text
}

func cacheAbuses() {
	CacheDirContent("profanity/data")
}

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

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

