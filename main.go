package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	under := kingpin.Flag("under", "Only show coverage lines under this percentage.").Default("100").Float32()
	dir := kingpin.Arg("package", "Package/directory to run coverage on.").String()
	kingpin.Parse()
	w, err := ioutil.TempFile(".", "gocoverlint-")
	kingpin.FatalIfError(err, "")
	defer os.Remove(w.Name())

	// (go test -coverprofile=t.cov 2> /dev/null; go tool cover -func=t.cov) | tr -d '%' | awk '/kingpin/ && $NF < 60 {print $1, $NF}'
	gotest := exec.Command("go", "test", "-coverprofile", w.Name(), *dir)
	gotest.Stderr = os.Stderr
	err = gotest.Run()
	kingpin.FatalIfError(err, "")

	cover := exec.Command("go", "tool", "cover", "-func", w.Name())
	cover.Stderr = os.Stderr
	out, err := cover.StdoutPipe()
	kingpin.FatalIfError(err, "")
	err = cover.Start()
	kingpin.FatalIfError(err, "")
	r := bufio.NewReader(out)
	for {
		bytes, _, e := r.ReadLine()
		if e != nil {
			break
		}
		line := string(bytes)
		if !strings.Contains(line, ".go:") {
			continue
		}
		fields := strings.Fields(line)
		coverage, e := strconv.ParseFloat(strings.TrimRight(fields[2], "%"), 32)
		kingpin.FatalIfError(e, "")
		if float32(coverage) <= *under {
			fmt.Println(line)
		}
	}
	err = cover.Wait()
	kingpin.FatalIfError(err, "")
}
