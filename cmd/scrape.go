package ace

import (
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/antchfx/htmlquery"
	"golang.org/x/net/html"
)

func ParseChallenges() []Challenge {
	challenges := make([]Challenge, 0, 600)
	source := getSource("https://www.techiedelight.com/data-structures-and-algorithms-problems/")
	challElements, err := htmlquery.QueryAll(source, "//div[@class='post-problems']//ol/li")
	if err != nil {
		log.Fatalf("Couldn't parse response Body, %s", err)
	}
	for i, el := range challElements {
		if err != nil {
			log.Fatalf("Couldn't parse response Body, %s", err)
		}
		challengeTitle := htmlquery.FindOne(el, "/a/text()")
		challengeUrl := htmlquery.SelectAttr(htmlquery.FindOne(el, "/a"), "href")
		source = getSource(challengeUrl)
		challengeDifficulty := htmlquery.FindOne(el, "/span/span/text()")
		var challengeTags []string
		for _, node := range htmlquery.Find(el, "/*[self::category or self::tag or self::lists]/text()") {
			challengeTags = slices.Concat(challengeTags, processTags(node.Data))
		}
		challengeDesc := htmlquery.Find(source, "//div[@class='post-content']/p[text()='For example,']/preceding-sibling::*/text()")
		var description string
		for _, t := range challengeDesc {
			description += t.Data
		}
		if i == 10 {
			break
		}
		challenges = append(challenges, Challenge{Description: description, Url: challengeUrl, Title: challengeTitle.Data, Difficulty: challengeDifficulty.Data, tags: challengeTags})
	}
	if err != nil {
		log.Fatalf("Couldn't parse response Body, %s", err)
	}
	return challenges
}

func processTags(tags string) []string {
	tags = strings.Replace(tags, ",", "", -1)
	tags = strings.ToLower(tags)
	return strings.Split(tags, " ")
}

func getSource(url string) *html.Node {
	cnx := http.Client{}
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
	resp, err := cnx.Do(req)
	if err != nil {
		log.Fatalf("Call to url %s failed \n %s", url, err)
	}
	defer resp.Body.Close()
	source, err := htmlquery.Parse(resp.Body)
	if err != nil {
		log.Fatalf("Couldn't parse response Body from %s \n %s", url, err)
	}
	return source
}
