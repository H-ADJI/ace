package ace

import (
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"
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
		number, err := strconv.Atoi(strings.Split(htmlquery.FindOne(el, "/text()").Data, ".")[0])
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
		if i == 50 {
			break
		}
		challenges = append(challenges, Challenge{id: number, description: description, url: challengeUrl, title: challengeTitle.Data, difficulty: challengeDifficulty.Data, tags: challengeTags})
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
	fmt.Println("visiting", url)
	resp, err := http.Get(url)
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
