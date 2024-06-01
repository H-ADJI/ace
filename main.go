package main

import (
	"fmt"

	ace "github.com/H-ADJI/ace/cmd"
)

func main() {
	challenges := ace.ParseChallenges()
	for i := range 10 {
		fmt.Println(challenges[i])
	}
}
