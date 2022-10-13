package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
)

func getHtmlPage(webPage string) (string, error) {
	resp, err := http.Get(webPage)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func getWebPageData(url string) string {
	data, _ := getHtmlPage(url)
	return data
}

func squares(c chan string, pages []string) {
	for i := 0; i <= len(pages)-1; i++ {
		x := countWordInPage(string(pages[i]))
		c <- strconv.Itoa(x)
	}
	close(c)
}

func countWordInPage(page string) int {
	data := getWebPageData(page)
	wordForCount := "Go"
	pattern := `[[:space:]|\W]` + wordForCount + `[[:space:]|\W]`
	regForAll := regexp.MustCompile(pattern)
	allTargetWords := regForAll.FindAllString(data, -1)
	count := len(allTargetWords)
	return count
}

func main() {
	pages := []string{"https://go.dev/", "https://go.dev/tour/moretypes/7", "https://go.dev/", "https://go.dev/tour/moretypes/7", "https://go.dev/tour/moretypes/7", "https://go.dev/", "https://go.dev/tour/moretypes/7", "https://go.dev/"}
	stack := 0

	chunkSize := 5
	result := make([][]string, 0)
	var first, last int
	for i := 0; i < len(pages)/chunkSize+1; i++ {
		first = i * chunkSize
		last = i*chunkSize + chunkSize
		if last > len(pages) {
			last = len(pages)
		}
		if first == last {
			break
		}

		result = append(result, pages[first:last])
	}
	for _, res := range result {
		c := make(chan string, 5)
		go squares(c, res)
		for val := range c {
			result, _ := strconv.Atoi(val)
			stack += result
		}

	}
	fmt.Println(stack)
}
