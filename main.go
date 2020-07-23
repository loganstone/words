package main

import (
	"container/heap"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"

	"jaytaylor.com/html2text"

	"github.com/loganstone/words/list"
)

const (
	pattern = "\\"
)

var isWord = regexp.MustCompile(`^[[:alpha:]]+$`).MatchString

// Words .
type Words []*word

func (w Words) Len() int {
	return len(w)
}

func (w Words) Less(i, j int) bool {
	if w[i].num == w[j].num {
		return w[i].txt < w[j].txt
	}
	return w[i].num < w[j].num
}

func (w Words) Swap(i, j int) {
	w[i], w[j] = w[j], w[i]
}

// Push .
func (w *Words) Push(x interface{}) {
	*w = append(*w, x.(*word))
}

// Pop .
func (w *Words) Pop() interface{} {
	old := *w
	n := len(old)
	element := old[n-1]
	*w = old[0 : n-1]
	return element
}

type word struct {
	txt string
	num int
}

var wg sync.WaitGroup

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage: %s <html page url>\n", os.Args[0])
		return
	}

	url := os.Args[1]

	res, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return
	}

	page, err := html2text.FromReader(res.Body, html2text.Options{})
	if err != nil {
		panic(err)
	}

	c := make(chan string)

	for _, line := range strings.Split(page, "\n") {
		data := strings.Split(line, " ")
		wg.Add(1)
		go func(data []string) {
			defer wg.Done()
			for _, datum := range data {
				w := strings.ToLower(datum)
				if !isWord(w) {
					continue
				}
				// TODO(logan): to option
				if list.IsExclude(w) {
					continue
				}

				c <- w
			}
		}(data)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	data := make(map[string]int)

	for w := range c {
		if _, ok := data[w]; ok == true {
			data[w]++
		} else {
			data[w] = 1
		}
	}

	words := &Words{}
	heap.Init(words)

	for txt, number := range data {
		heap.Push(words, &word{txt, number})
	}

	for words.Len() > 0 {
		w, ok := heap.Pop(words).(*word)
		if ok {
			fmt.Printf("%s, %d\n", w.txt, w.num)
		}
	}
}
