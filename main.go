package main

import (
	"fmt"

	ace "github.com/H-ADJI/ace/cmd"
)

func main() {
	challenges := ace.GetChallenges()
	for i := range 10 {
		fmt.Println(challenges[i])
	}
}
