package main

import (
	"bufio"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const WorkersNumber = 10
const ButchSize = 50

func main() {
	var cfg Config
	var filter Filter
	parseFlags(&cfg, &filter)

	//create ctx with timout equal to maxRunTime flag
	ctx, cancelFunc := context.WithTimeout(
		context.Background(),
		time.Duration(cfg.maxRunTime)*time.Second,
	)

	//fold existed ctx by notify context ( close when SIGINT/SIGTERM)
	ctx, _ = signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)

	//cancel only first context all another closed after this
	defer cancelFunc()

	//create filter funcs
	filterFuncs := newFilterFuncs(filter)

	file, err := os.Open(cfg.filePath)
	if err != nil {
		log.Fatal(err)
	}
	//close file when we and main
	defer file.Close()

	rawContents, err := gzip.NewReader(file)
	if err != nil {
		log.Fatal(err)
	}

	//run process of reading in workers
	err = processBulk(ctx, rawContents, NewJob(filterFuncs))
	if err != nil {
		log.Fatal(err)
	}

	if ctx.Err() != nil {
		fmt.Printf("interrupt with context error : %v", ctx.Err())
	}
}

func processBulk(ctx context.Context, r io.Reader, job func([]string) error) error {
	errs := make(chan error)
	input := make(chan []string)
	var wg sync.WaitGroup
	//create worker pool
	for i := 0; i < WorkersNumber; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				//wait for message or context close
				select {
				case bulk := <-input:
					if err := job(bulk); err != nil {
						errs <- err
						return
					}
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	//start reader
	go func() {
		Reader(ctx, r, input, errs)
	}()

	//close error chan when workers end
	go func() {
		wg.Wait()
		close(errs)
	}()

	//read from err chan and return err if happen
	var err error
	for e := range errs {
		if err == nil && e != nil {
			err = e
		}
	}
	return err
}

//Reader read from io.Reader by butches and send to workers
func Reader(ctx context.Context, r io.Reader, input chan []string, errs chan error) {
	scanner := bufio.NewScanner(r)
	bulk := make([]string, ButchSize)
	i := 0
	defer close(input)

	scanner.Scan() // read first line todo do something with column names
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return
		default:
			text := scanner.Text()
			bulk[i] = text
			i++
			if i == ButchSize {
				copied := make([]string, ButchSize, ButchSize)
				copy(copied, bulk)
				i = 0
				input <- copied
			}
		}
	}
	if i > 0 {
		input <- bulk[:i]
	}

	if err := scanner.Err(); err != nil {
		errs <- err
	}
}
