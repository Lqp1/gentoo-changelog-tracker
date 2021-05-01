package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/mmcdole/gofeed"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

func lookupAtomEquery(search string) string {
	atom_reg := regexp.MustCompile(`^([a-zA-Z_\-0-9]*/[a-zA-Z_\-0-9]*)-.*`)
	// TODO : Using equery is rather simple but we can't find packages that are not installed...
	//        Better use emerge something like that :
	// emerge --quiet --fuzzy-search n --search "%^ATOM$"
	cmd := exec.Command("equery", "-qC", "list", search)
	out, err := cmd.Output()
	if err != nil {
		log.Fatal("Atom not found on system; equery output: ", err)
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
func lookupAtomEix(search string) string {
	cmd := exec.Command("eix", "-#n", search)
	out, err := cmd.Output()
	if err != nil {
		log.Fatal("Atom not found on system; eix output: ", err)
	}
	atom := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	if len(atom) == 0 {
		log.Fatal("Could not match atom from eix output : " + string(out))
	}
	if len(atom) > 1 {
		fmt.Printf("Found atoms : %v\n", atom)
		fmt.Println("Several atoms were found. First one will be used")
	} else {
		fmt.Printf("Found atom : %s\n", atom[0])
	}
	return atom[0]
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
	blue := color.New(color.FgBlue).SprintFunc()
	header := fmt.Sprintf("[%d/%d]", index, length)
	time, _ := time.Parse(time.RFC3339, entry.Published)
	return fmt.Sprintf("%s %s on %s\n%s-----",
		red(header), green(entry.Title), blue(time), extractContent(entry.Content))
}

func formatDiff(commitID string) string {
	blue := color.New(color.FgBlue).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	var out string
	var inHeader bool = true
	uri := fmt.Sprintf("https://gitweb.gentoo.org/repo/gentoo.git/patch/?id=%s", commitID)
	diff, err := http.Get(uri)
	if err != nil {
		log.Fatal("Error while getting commit info from gitweb: ", err)
	}
	defer diff.Body.Close()
	rd := bufio.NewReader(diff.Body)
	// Only print maximum of 100 lines to avoid huge commits being printed
	// to stdout
	for i := 0; i < 100; i++ {
		str, err := rd.ReadString('\n')
		if err != nil {
			break
		}
		if strings.HasPrefix(str, "---") {
			inHeader = false
		}
		if inHeader == true {
			out += blue(str)
		} else {
			if strings.HasPrefix(str, "+++") || strings.HasPrefix(str, "---") || strings.HasPrefix(str, "@@") {
				out += blue(str)
			} else if strings.HasPrefix(str, "+") {
				out += green(str)
			} else if strings.HasPrefix(str, "-") {
				out += red(str)
			} else {
				out += str
			}
		}
	}
	return out
}

func fillCommit(queue chan struct {
	int
	string
}, index int, commit string) {
	queue <- struct {
		int
		string
	}{index, formatDiff(commit)}
}

func main() {
	limit := flag.Int("limit", 10, "How many entries to fetch")
	full := flag.Bool("full", false, "Print patch instead of commit summary")
	flag.Parse()

	if len(flag.Args()) != 1 {
		log.Fatal("Please provide exactly one atom")
	}
	if *limit > 10 {
		log.Fatal("Currently we cannot retrieve more than 10 entries")
	}

	atom := lookupAtomEix(flag.Args()[0])

	fp := gofeed.NewParser()
	uri := fmt.Sprintf("https://gitweb.gentoo.org/repo/gentoo.git/atom/%s?h=master", atom)
	feed, err := fp.ParseURL(uri)
	if err != nil {
		log.Fatal("Error while getting log info from gitweb: ", err)
	}
	if len(feed.Items) < *limit {
		*limit = len(feed.Items)
	}

	if *limit == 0 {
		fmt.Println("No entry to print")
	} else {
		if !*full {
			for i := 0; i < *limit; i++ {
				fmt.Println(formatEntry(feed.Items[i], i+1, *limit))
			}
		} else {
			queue := make(chan struct {
				int
				string
			})
			output := make([]string, *limit)
			for i := 0; i < *limit; i++ {
				go fillCommit(queue, i, feed.Items[i].GUID)
			}
			for i := 0; i < *limit; i++ {
				diff := <-queue
				output[diff.int] = diff.string
			}
			fmt.Println(strings.Join(output, "\n"))
		}
	}

}
