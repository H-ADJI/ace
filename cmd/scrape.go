package ace

import (
	"log"
	"net/http"
	"slices"
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
		challengeDifficulty := htmlquery.FindOne(el, "/span/span/text()")
		var challengeTags []string
		for _, node := range htmlquery.Find(el, "/*[self::category or self::tag or self::lists]/text()") {
			challengeTags = slices.Concat(challengeTags, processTags(node.Data))
		}
		challenges = append(challenges, Challenge{id: number, url: challengeUrl, title: challengeTitle.Data, difficulty: challengeDifficulty.Data, tags: challengeTags})
	}
	if err != nil {
		log.Fatalf("Couldn't parse response Body, %s", err)
	}
	return challenges
}

func processTags(tags string) []string {

	// remove punc
	tags = strings.Replace(tags, ",", "", -1)
	// lower case
	tags = strings.ToLower(tags)
	// split
	return strings.Split(tags, " ")
}
