package ace

import (
	"database/sql"
	"log"
	"net/http"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

const workerPoolSize = 3

func crawlChallenges(db *sql.DB) {
	c := http.Client{}
	source := getSource(c, "https://www.techiedelight.com/data-structures-and-algorithms-problems/")
	challElements, err := htmlquery.QueryAll(source, "//div[@class='post-problems']//ol/li")
	pubChan := make(chan Challenge, 5)
	consumeChan := make(chan Challenge, 5)
	if err != nil {
		log.Fatalf("Couldn't parse response Body, %s", err)
	}
	for i := 0; i < workerPoolSize; i++ {
		go addDescription(pubChan, consumeChan, c)
	}
	go func() {
		for chall := range consumeChan {
			chall.InsertIntoDB(db)
		}
	}()
	for _, el := range challElements {
		if err != nil {
			log.Fatalf("Couldn't parse response Body, %s", err)
		}
		challengeTitle := htmlquery.FindOne(el, "/a/text()")
		challengeUrl := htmlquery.SelectAttr(htmlquery.FindOne(el, "/a"), "href")
		challengeDifficulty := htmlquery.FindOne(el, "/span/span/text()")
		var challengeTags strings.Builder
		for _, node := range htmlquery.Find(el, "/*[self::category or self::tag or self::lists]/text()") {
			challengeTags.WriteString(node.Data)
			challengeTags.WriteString(" ")
		}
		pubChan <- Challenge{
			Url:        challengeUrl,
			Title:      challengeTitle.Data,
			Difficulty: challengeDifficulty.Data,
			Tags:       challengeTags.String()}
	}
	close(pubChan)

}
func addDescription(receiveChan <-chan Challenge, sendChan chan<- Challenge, client http.Client) {
	for chall := range receiveChan {
		source := getSource(client, chall.Url)
		challengeDesc := htmlquery.Find(source, "//div[@class='post-content']/p[text()='For example,']/preceding-sibling::*/text()")
		var description strings.Builder
		for _, t := range challengeDesc {
			description.WriteString(t.Data)
			description.WriteString(" ")
		}
		chall.Description = description.String()
		sendChan <- chall
	}
	close(sendChan)
}
func getSource(client http.Client, url string) *html.Node {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Build request failed\n %s", err)
	}
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8")
	req.Header.Set("accept-language", "en-US,en;q=0.7")
	req.Header.Set("cache-control", "no-cache")
	req.Header.Set("pragma", "no-cache")
	req.Header.Set("priority", "u=0, i")
	req.Header.Set("referer", "https://www.techiedelight.com/data-structures-and-algorithms-problems/")
	req.Header.Set("sec-ch-ua", `"Brave";v="125", "Chromium";v="125", "Not.A/Brand";v="24"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Linux"`)
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "none")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("sec-gpc", "1")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/125.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Call to url %s failed \n %s", url, err)
	}
	defer resp.Body.Close()
	source, err := htmlquery.Parse(resp.Body)
	if err != nil {
		log.Fatalf("Couldn't parse response Body from %s \n %s", url, err)
	}
	log.Println("Got response for ...", url)
	return source
}
