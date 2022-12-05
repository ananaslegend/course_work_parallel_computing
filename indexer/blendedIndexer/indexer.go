package blendedIndexer

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"text/tabwriter"
)

func BuildIndex(path string, threadsCount int) (*InvertedIndex, error) {
	if threadsCount < 1 {
		return nil, errors.New("thread count should be grater then 0")
	}
	docs, err := os.ReadDir(path)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	jobs := make(chan string, len(docs))

	go func() {
		for _, doc := range docs {
			jobs <- doc.Name()
		}
		close(jobs)
	}()

	wg := sync.WaitGroup{}

	ii := newInvertedIndex()

	for i := 0; i < threadsCount; i++ {
		wg.Add(1)

		go func() {
			for doc := range jobs {
				file, err := os.Open(path + doc)
				if err != nil {
					log.Println(err)
					wg.Done()
					return
				}

				scanner := bufio.NewScanner(file)
				scanner.Split(bufio.ScanWords)

				for scanner.Scan() {
					word := strings.Trim(scanner.Text(), ".,/:';!@#$%&*()`~<>[]{}\n\r")

					ii.lock.Lock()
					ii.insert(word, doc)
					ii.lock.Unlock()
				}

				if err := file.Close(); err != nil {
					log.Println(err)
					wg.Done()
					return
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()

	return ii, nil
}

func (ii *InvertedIndex) PrintSingleWordEntering(word string) {
	em := ii.m[word]
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "\nFile Name\t Entering")
	for file, ent := range em {
		fmt.Fprintln(w, file, "\t", ent)
	}

	if err := w.Flush(); err != nil {
		log.Println(err)
		return
	}
}

func (ii *InvertedIndex) PrintIndexTable() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "\nWord\t [File Name : Entering]")

	for word, files := range ii.m {
		str := fmt.Sprintf("%s\t", word)
		for filename, number := range files {
			str += fmt.Sprintf(" [%s : %d]", filename, number)
		}
		fmt.Fprintln(w, str)
	}

	if err := w.Flush(); err != nil {
		log.Println(err)
		return
	}
}

type InvertedIndex struct {
	lock *sync.Mutex
	m    map[string]map[string]int
}

func (ii *InvertedIndex) insert(word string, fileName string) {
	if ii.m[word] == nil {
		ii.m[word] = make(map[string]int)
	}

	ii.m[word][fileName] += 1
}

func newInvertedIndex() *InvertedIndex {
	return &InvertedIndex{lock: &sync.Mutex{}, m: make(map[string]map[string]int)}
}
