package main

import (
	"io/ioutil"
	"os"
	"os/exec"

	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	dir := kingpin.Arg("package", "Package/directory to run coverage on.").String()
	kingpin.Parse()
	w, err := ioutil.TempFile(".", "gocoverlint-")
	kingpin.FatalIfError(err, "")
	defer os.Remove(w.Name())

	// (go test -coverprofile=t.cov 2> /dev/null; go tool cover -func=t.cov) | tr -d '%' | awk '/kingpin/ && $NF < 60 {print $1, $NF}'
	gotest := exec.Command("go", "test", "-coverprofile", w.Name(), *dir)
	gotest.Stdout = os.Stdout
	gotest.Stderr = os.Stderr
	err = gotest.Run()
	kingpin.FatalIfError(err, "")

	cover := exec.Command("go", "tool", "cover", "-func", w.Name())
	cover.Stdout = os.Stdout
	cover.Stderr = os.Stderr
	err = cover.Run()
	kingpin.FatalIfError(err, "")
}
