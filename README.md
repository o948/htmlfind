# htmlfind
Small library to find HTML nodes with CSS-like selectors.

```go
type Node struct {
        Attr     map[string]string
        Children []*Node
        // Has unexported fields.
}

func Parse(r io.Reader) (*Node, error)
func (root *Node) FindFirst(query string) *Node
func (root *Node) FindN(query string, n int) []*Node
func (root *Node) Tag() string
func (root *Node) Text() string
```

Query consists of multiple node selectors separated by whitespace. Node selector is:
```
[>] element [.class_name | #id_value] [:nth_child | :-nth_last_child]
```

## Example
Let's get the list of Go standard library packages from https://golang.org/pkg/:
```go
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/o948/htmlfind"
)

func main() {
	resp, err := http.Get("https://golang.org/pkg/")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	root, err := htmlfind.Parse(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	for _, a := range root.FindN("td.pkg-name a", -1) {
		fmt.Println(a.Attr["href"])
	}
}
```
