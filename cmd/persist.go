package ace

import "fmt"

type Challenge struct {
	title       string
	id          int
	url         string
	description string
	isSolved    bool
	tags        []string
	difficulty  string
}

func (chall Challenge) String() string {
	return fmt.Sprintf("Title : %s  |  Difficulty : %s ", chall.title, chall.difficulty)
}
