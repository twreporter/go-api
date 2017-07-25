package main

import (
	"bufio"
	"fmt"
	"go/build"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	p, _ := build.Default.Import("twreporter.org/go-api", "", build.FindOnly)

	fname := filepath.Join(p.Dir, "CHANGELOG.md")
	file, err := os.Open(fname)
	if err != nil {
		fmt.Println("Cannot open CHANGELOG.md")
		return
	}

	defer file.Close()
	reader := bufio.NewReader(file)
	line, _, err1 := reader.ReadLine()
	if err1 != nil {
		fmt.Println("Cannot read file")
		fmt.Println(err1.Error())
		return
	}
	ver := strings.Replace(string(line), "#", "", -1)
	ver = strings.Replace(ver, " ", "", -1)

	fmt.Print(ver)
}
