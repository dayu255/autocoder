package fetch

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Test struct {
	Num     int
	Type      string
	Content string
}

func titleParse(title string) (num int, io string) {
	token := strings.Split(title, " ")

	switch token[1] {
	case "Input":
		io = "input"
	case "Output":
		io = "output"
	default:
		io = "unknown"
	}

	num, _ = strconv.Atoi(token[2])
	return
}

func Fetch(url string) []Test {
	client := http.Client{}
	res, _ := client.Get(url)

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		panic(err)
	}

	tests := make([]Test, 0, 12);

	doc.Find(".part").Each(func(i int, s *goquery.Selection) {
		title := s.Find("h3").Text()

		if strings.Contains(title, "Sample Input") || strings.Contains(title, "Sample Output") {
			content := s.Find("pre").Text()

			n, t := titleParse(title)

			tests = append(tests, Test{n ,t, content})
		}
	})

	return tests
}
