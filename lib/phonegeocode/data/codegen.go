package main

import (
	"bufio"
	"embed"
	"encoding/csv"
	"fmt"
	"io"
	"log"
)

//go:embed prefixes.csv
var prefixes embed.FS

func main() {
	f, err := prefixes.Open("prefixes.csv")
	if err != nil {
		log.Fatal(err)
	}
	in := bufio.NewReader(f)
	reader := csv.NewReader(in)

	fmt.Printf("package phonegeocode\n\n")
	fmt.Printf("import (\n\tgotrie \"github.com/tchap/go-patricia/patricia\"\n)\n\n")

	fmt.Printf("func initPrefixes() *gotrie.Trie {\n")
	fmt.Printf("\tprefixes := gotrie.NewTrie()\n\n")

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("	prefixes.Insert(gotrie.Prefix(\"%s\"), \"%s\")\n", row[0], row[1])
	}
	fmt.Printf("\n\treturn prefixes\n}\n")
}
