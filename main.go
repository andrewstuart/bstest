package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

var pkgRE = regexp.MustCompile(`package (?P<package>\w+)`)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("please pass percentage float64 as first argument")
	}

	pct, err := strconv.ParseFloat(os.Args[1], 64)
	if err != nil {
		log.Fatal(err)
	}

	if pct < 0.0 || pct >= 100.0 {
		log.Fatal("please enter a percentage > 0 and < 100")
	}

	pct /= 100.0

	fs, err := ioutil.ReadDir(".")
	if err != nil {
		log.Fatal(err)
	}

	lines := 0
	wd, err := os.Getwd()
	if err != nil {
		wd = "main"
	}
	pkg := path.Base(wd)

	for _, fi := range fs {
		n := fi.Name()
		if !fi.IsDir() && strings.HasSuffix(n, ".go") && !strings.HasSuffix(n, "_test.go") && n != "bs.go" && n != "bs_test.go" {
			f, err := ioutil.ReadFile(fi.Name())
			if err != nil {
				log.Fatal(err)
			}

			lines += bytes.Count(f, []byte{'\n'})
			for _, submatches := range pkgRE.FindAllSubmatchIndex(f, -1) {
				bs := pkgRE.Expand([]byte{}, []byte("$package"), f, submatches)
				if len(bs) > 0 {
					pkg = string(bs)
				}
			}
		}
	}

	tgt := 1.0 - pct

	gen := float64(lines) / tgt

	f, err := os.OpenFile("bs.go", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0664)
	if err != nil {
		log.Fatal(err)
	}
	bw := bufio.NewWriter(f)

	fmt.Fprintf(f, "package %s\nfunc bs() {\n\tbs := 1\n", pkg)
	for i := 0; i < int(gen); i++ {
		bw.Write([]byte("\tbs = 1\n"))
	}
	bw.Flush()
	fmt.Fprintf(f, "\n_ = bs\n}\n")
	f.Close()

	f, err = os.OpenFile("bs_test.go", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0664)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Fprintf(f, `package %s
import "testing"
func TestBS(t *testing.T) {
	bs()
}
`, pkg)
}
