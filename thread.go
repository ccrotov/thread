package thread

import (
	"bufio"
	"errors"
	"os"
	"sync"
)

var (
	isWork bool = false
)

func FileEnum(fname string, threads int, cback func(ln string)) error {
	defer Stop()
	isWork = true

	var lines []string

	chn := make(chan string, threads)

	file, err := os.Open(fname)
	if err != nil {
		return err
	}

	scaner := bufio.NewScanner(file)

	for scaner.Scan() {
		lines = append(lines, scaner.Text())
	}

	if len(lines) == 0 {
		return errors.New("file is empty " + fname)
	}

	err = file.Close()
	if err != nil {
		return err
	}

	if scaner.Err() != nil {
		return scaner.Err()
	}

	var wg sync.WaitGroup
	for i := 0; i < cap(chn); i++ {
		go worker(chn, &wg, cback)
	}

	for _, line := range lines {
		wg.Add(1)
		chn <- line
	}

	wg.Wait()
	close(chn)
	return nil
}

func Stop() {
	isWork = false
}

func worker(ch chan string, wg *sync.WaitGroup, cback func(ln string)) {
	for p := range ch {

		if !isWork {
			wg.Done()
			continue
		}

		cback(p)
		wg.Done()
	}
}
