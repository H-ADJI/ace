package main

import (
	"fmt"
	"net/http"

	"github.com/antchfx/htmlquery"
)

func main() {
	resp, err := http.Get("https://www.techiedelight.com/data-structures-and-algorithms-problems/")
	if err != nil {
		fmt.Println("something went wrong")
		return
	}
	defer resp.Body.Close()
	source, err := htmlquery.Parse(resp.Body)
	if err != nil {
		fmt.Println("couldnt parse resp body")
		return
	}
	elmnts, err := htmlquery.Query(source, "//div")
	if err != nil {
		fmt.Println("couldnt query source")

	}
	elmnts = htmlquery.FindOne(elmnts, "//h1/text()")
	fmt.Println(elmnts.Data)
}
