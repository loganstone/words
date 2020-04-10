package main

import (
	"container/heap"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"jaytaylor.com/html2text"
)

const (
	pattern      = "\\"
	prepositions = `aboard
about
above
absent
across
after
against
along
alongside
amid
among
amongst
anti
around
as
at
before
behind
below
beneath
beside
besides
between
beyond
but
by
circa
concerning
considering
despite
down
during
except
excepting
excluding
failing
following
for
from
given
in
inside
into
like
minus
near
of
off
on
onto
opposite
outside
over
past
per
plus
regarding
round
save
since
than
through
to
toward
towards
under
underneath
unlike
until
up
upon
versus
via
with
within
without
worth`
)

var isWord = regexp.MustCompile(`^[[:alpha:]]+$`).MatchString
var mapByPreposition map[string]bool

// Words .
type Words []*word

func (w Words) Len() int {
	return len(w)
}

func (w Words) Less(i, j int) bool {
	return w[i].num > w[j].num
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

func init() {
	tmp := make(map[string]bool)
	data := strings.Split(prepositions, "\n")
	for _, datum := range data {
		tmp[datum] = true
	}
	mapByPreposition = tmp
}

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

	data := make(map[string]int)
	lines := strings.Split(page, "\n")
	for _, line := range lines {
		for _, datum := range strings.Split(line, " ") {
			w := strings.ToLower(datum)
			if !isWord(w) {
				continue
			}
			// TODO(logan): to option
			if _, ok := mapByPreposition[w]; ok {
				continue
			}
			if _, ok := data[w]; ok {
				data[w]++
			} else {
				data[w] = 1
			}
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
