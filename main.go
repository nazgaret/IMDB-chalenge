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
	ctx, cancelFunc := context.WithTimeout(
		context.Background(),
		time.Duration(cfg.maxRunTime)*time.Second,
	)
	ctx, _ = signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer cancelFunc()

	filterFuncs := newFilterFuncs(filter)

	file, err := os.Open(cfg.filePath)
	if err != nil {
		log.Fatal(err)
	}

	rawContents, err := gzip.NewReader(file)
	if err != nil {
		log.Fatal(err)
	}

	err = processBulk(ctx, rawContents, NewJob(filterFuncs))
	if err != nil {
		log.Fatal(err)
	}
}

func processBulk(ctx context.Context, r io.Reader, job func([]string) error) error {
	errs := make(chan error)
	input := make(chan []string)
	var wg sync.WaitGroup
	for i := 0; i < WorkersNumber; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
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

	go func() {
		scanner := bufio.NewScanner(r)
		l := ButchSize
		bulk := make([]string, l)
		i := 0
		defer close(input)

		scanner.Scan() // read first line todo do something
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				fmt.Println(ctx.Err())
				return
			default:
				text := scanner.Text()
				bulk[i] = text
				i++
				if i == l {
					copied := make([]string, l, l)
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
	}()
	go func() {
		wg.Wait()
		close(errs)
	}()
	var err error
	for e := range errs {
		if err == nil && e != nil {
			err = e
		}
	}
	return err
}
