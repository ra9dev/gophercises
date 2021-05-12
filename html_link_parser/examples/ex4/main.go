package main

import (
	"fmt"
	"strings"

	link "github.com/KubaiDoLove/gophercises_html_link_parser"
)

var exampleHTML = `
<html>
<body>
  <a href="/dog-cat">dog cat <!-- commented text SHOULD NOT be included! --></a>
</body>
</html>
`

func main() {
	r := strings.NewReader(exampleHTML)

	links, err := link.Parse(r)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", links)
}
