package ace

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
)

func GetChallenges() []Challenge {
	challenges := make([]Challenge, 0, 600)
	resp, err := http.Get("https://www.techiedelight.com/data-structures-and-algorithms-problems/")
	if err != nil {
		log.Fatalf("Call to website failed, %s", err)
	}
	defer resp.Body.Close()
	source, err := htmlquery.Parse(resp.Body)
	if err != nil {
		log.Fatalf("Couldn't parse response Body, %s", err)
	}
	challElements, err := htmlquery.QueryAll(source, "//div[@class='post-problems']//ol/li")
	if err != nil {
		log.Fatalf("Couldn't parse response Body, %s", err)
	}
	for _, el := range challElements {
		number, err := strconv.Atoi(strings.Split(htmlquery.FindOne(el, "/text()").Data, ".")[0])
		if err != nil {
			log.Fatalf("Couldn't parse response Body, %s", err)
		}
		challengeTitle := htmlquery.FindOne(el, "/a/text()")
		challengeUrl := htmlquery.SelectAttr(htmlquery.FindOne(el, "/a"), "href")
		// challengeDifficulty := htmlquery.FindOne(el, "/span/span/text()")
		// challengeTags := htmlquery.FindOne(el, "/*[self::category or self::tag or self::lists]/text()")
		challenges = append(challenges, Challenge{id: number, url: challengeUrl, title: challengeTitle.Data})
	}
	if err != nil {
		log.Fatalf("Couldn't parse response Body, %s", err)
	}
	return challenges
}
