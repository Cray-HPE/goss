package goss

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/aelsabbahy/goss/resource"
	"github.com/aelsabbahy/goss/system"
	"github.com/codegangsta/cli"
)

func Run(specFile string, c *cli.Context) {
	sys := system.New(c)

	// handle stdin
	var fh *os.File
	var err error
	var path string
	if hasStdin() {
		fh = os.Stdin
	} else {
		path = filepath.Dir(specFile)
		fh, err = os.Open(specFile)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}
	data, err := ioutil.ReadAll(fh)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	configJSON := mergeJSONData(ReadJSONData(data), 0, path)

	out := make(chan resource.TestResult)

	in := make(chan resource.Resource)

	go func() {
		for _, t := range configJSON.Resources() {
			in <- t
		}
		close(in)
	}()

	if os.Getenv("GOMAXPROCS") == "" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	gomaxprocs := runtime.GOMAXPROCS(-1)
	workerCount := gomaxprocs * 5
	if workerCount > 50 {
		workerCount = 50
	}
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for f := range in {
				for _, r := range f.Validate(sys) {
					out <- r
				}
			}

		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	testCount := 0
	var failed []resource.TestResult
	for testResult := range out {
		//fmt.Printf("%v: %s.\n", testResult.Duration, testResult.Desc)
		if testResult.Result {
			fmt.Printf(".")
			testCount++
		} else {
			fmt.Printf("F")
			failed = append(failed, testResult)
			testCount++
		}
	}

	for _, testResult := range failed {
		fmt.Printf("\n%s\n", testResult.Desc)
	}

	fmt.Printf("\n\nCount: %d failed: %d\n", testCount, len(failed))
	if len(failed) > 0 {
		os.Exit(1)
	}
}

func hasStdin() bool {
	if fi, err := os.Stdin.Stat(); err == nil {
		mode := fi.Mode()
		if (mode&os.ModeNamedPipe != 0) || mode.IsRegular() {
			return true
		}
	}
	return false
}