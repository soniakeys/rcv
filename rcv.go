// Rcv, a tool to add or update a test coverage section in a readme.
//
// Run rcv in a directory with a readme.md.  If go test -cover finds any
// test coverage, rcv updates the readme.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"time"
)

const readme = "readme.md"
const tcHdr = "###Test coverage"

func main() {
	// load existing readme
	rm, err := ioutil.ReadFile(readme)
	if err != nil {
		log.Fatal(err)
	}

	// clip existing test coverage
	if x := regexp.MustCompile(`\n` + tcHdr + `\n`).FindIndex(rm); x != nil {
		rm = rm[:x[0]]
	}

	// get coverage output
	c, err := exec.Command("go", "test", "-cover", "./...").Output()
	if err != nil {
		log.Fatal(err)
	}

	// scrape the parts we want
	d, _ := os.Getwd()
	_, wd := path.Split(d)
	r := regexp.MustCompile(`ok .*(` + wd + `.*)\t\d.*coverage.*?(\d+\.\d+%)`)
	x := r.FindAllSubmatch(c, -1)
	if x == nil {
		log.Fatal("no coverage data")
	}

	// overwrite readme
	f, err := os.Create(readme)
	if err != nil {
		log.Fatal(err)
	}
	f.Write(rm)
	fmt.Fprintln(f, "\n"+tcHdr)
	fmt.Fprintln(f, time.Now().Format("2 Jan 2006"))
	fmt.Fprintln(f, "```")
	max := 0
	for _, s := range x {
		if len(s[1]) > max {
			max = len(s[1])
		}
	}
	for _, s := range x {
		fmt.Fprintf(f, "%-*s  %s\n", max, string(s[1]), string(s[2]))
	}
	fmt.Fprintln(f, "```")
	f.Close()
}
