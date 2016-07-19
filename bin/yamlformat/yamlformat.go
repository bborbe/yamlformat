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
	"runtime"
)

const (
	PARAMETER_LOGLEVEL = "loglevel"
	PARAMETER_PATH     = "path"
	PARAMETER_WRITE    = "write"
)

var (
	logger      = log.DefaultLogger
	logLevelPtr = flag.String(PARAMETER_LOGLEVEL, log.LogLevelToString(log.ERROR), log.FLAG_USAGE)
	pathPtr     = flag.String(PARAMETER_PATH, "", "path")
	writePtr    = flag.Bool(PARAMETER_WRITE, false, "write")
)

func main() {
	defer logger.Close()
	flag.Parse()

	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	runtime.GOMAXPROCS(runtime.NumCPU())

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
	logger.Debugf("format %s and write %v", path, write)
	if len(path) == 0 {
		fmt.Fprintf(writer, "parameter %s missing\n", PARAMETER_PATH)
		return nil
	}
	if path, err = io_util.NormalizePath(path); err != nil {
		logger.Warnf("normalize path: %s failed: %v", path, err)
		return err
	}
	source, err := ioutil.ReadFile(path)
	if err != nil {
		logger.Warnf("read file %s failed: %v", path, err)
		return err
	}
	var data interface{}
	err = yaml.Unmarshal(source, &data)
	if err != nil {
		logger.Warnf("unmarshal %s failed: %v", path, err)
		return err
	}
	content, err := yaml.Marshal(data)
	if err != nil {
		logger.Warnf("marshal failed: %v", err)
		return err
	}
	if write {
		logger.Debug("write file")
		fileInfo, err := readMode(path)
		if err != nil {
			logger.Warnf("get fileinfo failed: %v", err)
			return err
		}
		logger.Debugf("write yaml %s", path)
		return ioutil.WriteFile(path, content, fileInfo.Mode())
	} else {
		logger.Debugf("print content")
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
