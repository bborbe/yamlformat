package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"runtime"

	io_util "github.com/bborbe/io/util"
	"github.com/golang/glog"
	"gopkg.in/yaml.v2"
)

const (
	parameterPath  = "path"
	parameterWrite = "write"
)

var (
	pathPtr  = flag.String(parameterPath, "", "path")
	writePtr = flag.Bool(parameterWrite, false, "write")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	writer := os.Stdout

	err := do(writer, *pathPtr, *writePtr)
	if err != nil {
		glog.Exitf("format %v failed: %v", *pathPtr, err)
	}
	glog.V(2).Infof("format %v success", *pathPtr)
}

func do(writer io.Writer, path string, write bool) error {
	var err error
	glog.V(2).Infof("format %s and write %v", path, write)
	if len(path) == 0 {
		fmt.Fprintf(writer, "parameter %s missing\n", parameterPath)
		return nil
	}
	if path, err = io_util.NormalizePath(path); err != nil {
		glog.Warningf("normalize path: %s failed: %v", path, err)
		return err
	}
	source, err := ioutil.ReadFile(path)
	if err != nil {
		glog.Warningf("read file %s failed: %v", path, err)
		return err
	}
	var data interface{}
	err = yaml.Unmarshal(source, &data)
	if err != nil {
		glog.Warningf("unmarshal %s failed: %v", path, err)
		return err
	}
	content, err := yaml.Marshal(data)
	if err != nil {
		glog.Warningf("marshal failed: %v", err)
		return err
	}
	if write {
		glog.V(2).Info("write file")
		fileInfo, err := readMode(path)
		if err != nil {
			glog.Warningf("get fileinfo failed: %v", err)
			return err
		}
		glog.V(2).Infof("write yaml %s", path)
		return ioutil.WriteFile(path, content, fileInfo.Mode())
	} else {
		glog.V(2).Infof("print content")
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
