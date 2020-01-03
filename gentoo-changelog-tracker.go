package main

import (
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/mmcdole/gofeed"
	"golang.org/x/net/html"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

func lookupAtom(search string) string {
	atom_reg := regexp.MustCompile(`^(.*/.*?)-.*`)
	cmd := exec.Command("equery", "-qC", "list", search)
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	atom := atom_reg.FindStringSubmatch(string(out))
	if len(atom) <= 1 {
		log.Fatal("Could not match atom from equery output : " + string(out))
	}
	if len(atom) > 2 {
		fmt.Printf("Found atoms : %v\n", atom[1:])
		fmt.Println("Several atoms were found. First one will be used")
	} else {
		fmt.Printf("Found atom : %s\n", atom[1])
	}
	return atom[1]
}

func extractContent(src string) string {
	var crawler func(*html.Node)
	var out string
	doc, err := html.Parse(strings.NewReader(src))
	if err != nil {
		log.Fatal(err)
	}

	crawler = func(node *html.Node) {
		if node.Type == html.TextNode {
			out += html.UnescapeString(node.Data)
		}
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			crawler(child)
		}
	}
	crawler(doc)
	if out == "" {
		return html.UnescapeString(src)
	} else {
		return strings.TrimSuffix(out, "\n")
	}
}

func formatEntry(entry *gofeed.Item, index int, length int) string {
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	header := fmt.Sprintf("[%d/%d]", index, length)
	return fmt.Sprintf("%s %s\n%s-----",
		red(header), green(entry.Title), extractContent(entry.Content))
}

func main() {
	limit := flag.Int("limit", 10, "How many entries to fetch")
	flag.Parse()

	if len(flag.Args()) != 1 {
		log.Fatal("Please provide exactly one atom")
	}
	if *limit > 10 {
		log.Fatal("Currently we cannot retrieve more than 10 entries")
	}

	atom := lookupAtom(flag.Args()[0])

	fp := gofeed.NewParser()
	uri := fmt.Sprintf("https://gitweb.gentoo.org/repo/gentoo.git/atom/%s?h=master", atom)
	feed, err := fp.ParseURL(uri)
	if err != nil {
		log.Fatal(err)
	}
	if len(feed.Items) < *limit {
		*limit = len(feed.Items)
	}

	if *limit == 0 {
		fmt.Println("No entry to print")
	} else {
		for i := 0; i < *limit; i++ {
			fmt.Println(formatEntry(feed.Items[i], i+1, *limit))
		}
	}

}
