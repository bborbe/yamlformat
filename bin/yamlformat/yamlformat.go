package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	io_util "github.com/bborbe/io/util"
	"github.com/bborbe/log"
	"gopkg.in/yaml.v2"
)

var logger = log.DefaultLogger

const (
	PARAMETER_LOGLEVEL = "loglevel"
	PARAMETER_PATH     = "path"
	PARAMETER_WRITE    = "write"
)

func main() {
	defer logger.Close()
	logLevelPtr := flag.String(PARAMETER_LOGLEVEL, log.LogLevelToString(log.ERROR), log.FLAG_USAGE)
	pathPtr := flag.String(PARAMETER_PATH, "", "path")
	writePtr := flag.Bool(PARAMETER_WRITE, false, "write")
	flag.Parse()
	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	writer := os.Stdout

	err := do(writer, *pathPtr, *writePtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
	logger.Debug("done")
}

func do(writer io.Writer, path string, write bool) error {
	var err error
	if len(path) == 0 {
		fmt.Fprintf(writer, "parameter %s missing\n", PARAMETER_PATH)
		return nil
	}
	if path, err = io_util.NormalizePath(path); err != nil {
		return err
	}
	source, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	var data interface{}
	err = yaml.Unmarshal(source, &data)
	if err != nil {
		return err
	}
	content, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	if write {
		fileInfo, err := readMode(path)
		if err != nil {
			return err
		}
		logger.Debugf("write yaml %s", path)
		return ioutil.WriteFile(path, content, fileInfo.Mode())
	} else {
		fmt.Fprintf(writer, "%s\n", string(content))
		return nil
	}
}

func readMode(path string) (os.FileInfo, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return nil, err
	}
	return f.Stat()
}
