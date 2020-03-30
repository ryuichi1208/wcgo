package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"sync"
	"time"

	"github.com/golang/glog"
)

type Arg struct {
	targetFile string
	splitNum int
	maxThreads int
	buffersize int
}

var (
	arg Arg
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "%s\n", fmt.Sprintf("%s -f TARGETFILE [options] [glog options]", os.Args[0]))
		flag.PrintDefaults()
	}

	_ = flag.Set("stderrthreshold", "INFO")
	_ = flag.Set("v", "0")

	flag.StringVar(&arg.targetFile, "f", "", "(go-lc) Target File")
	flag.IntVar(&arg.splitNum, "s", 2, "(go-lc) Num of File split")
	flag.IntVar(&arg.maxThreads, "t", 2, "(go-lc) Max Num of Threads")
	flag.IntVar(&arg.buffersize, "b", 1024*1024, "(go-lc) Size of ReadBuffer")
}

func getFileSize(filename string) (int, error) {
	fh, err := os.OpenFile(filename, 0, 0)
	if err != nil {
		return 0, err
	}
	defer fh.Close()

	fileinfo, err := fh.Stat()
	if err != nil {
		return 0, err
	}

	return int(fileinfo.Size()), nil
}

func getNumOfLines(filename string, splitNum int, maxThreads int, buffersize int) (int, error) {
	fsize, err := getFileSize(filename)
	if err != nil {
		return 0, err
	}

	glog.V(1).Infof("FileSize   : %10d byte", fsize)
	glog.V(1).Infof("Read buffer: %10d byte", buffersize)
	glog.V(1).Infof("Max Threads: %d", maxThreads)
	glog.V(1).Infof("Split Num  : %d", splitNum)

	var readCountTotal int = int(math.Trunc(float64(fsize) / float64(buffersize)))

	if fsize-(readCountTotal*buffersize) > 0 {
		readCountTotal++
	}

	wg := &sync.WaitGroup{}

	jc := make(chan interface{}, maxThreads)
	defer close(jc)

	counterCh := make(chan int, maxThreads)

	resultCh := make(chan int)
	defer close(resultCh)

	go func(counterCh <-chan int) {
		cAll := 0
		for c := range counterCh {
			cAll += c

			glog.V(2).Infof("[receiver] receive: %d\n", c)
		}

		resultCh <- cAll
	}(counterCh)

	var byteOffset int64 = 0

	for i := 0; i < splitNum; i++ {
		eachReadCount := int(math.Trunc(float64(readCountTotal+i) / float64(splitNum)))
		jc <- true
		wg.Add(1)
		go countWorker(filename, eachReadCount, byteOffset, buffersize, wg, jc, counterCh)
		byteOffset += int64(eachReadCount * buffersize)
	}

	wg.Wait()
	close(counterCh)

	return <-resultCh, nil
}

func main() {
	flag.Parse()
	glog.V(1).Infof("Start")
	startTime := time.Now()
	numOfLines, _ := getNumOfLines(arg.targetFile, arg.splitNum, arg.maxThreads, arg.buffersize)
	glog.V(1).Infof("End(%s)", time.Since(startTime))
	fmt.Printf("%d\n", numOfLines)
}
