package ace

type Challenge struct {
	title       string
	id          int
	url         string
	description string
	isSolved    bool
	tags        []string
	difficulty  int8
}
