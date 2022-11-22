package blendedIndexer

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
)

func BuildIndex(path string, threadsCount int) (*InvertedIndex, error) {
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
				file, err := os.Open(doc)
				if err != nil {
					log.Println(err)
					wg.Done()
					return
				}

				scanner := bufio.NewScanner(file)
				scanner.Split(bufio.ScanWords)

				for scanner.Scan() {
					word := scanner.Text()

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
		wg.Wait()
	}
	return ii, nil
}

func (ii *InvertedIndex) FindWordEntering() {

}

func (ii *InvertedIndex) Print() {
	for word, files := range ii.m {
		str := fmt.Sprintf("%s: ", word)
		for filename, number := range files {
			str += fmt.Sprintf("[%s : %d]", filename, number)
		}
		fmt.Println(str)
	}
}

type InvertedIndex struct {
	lock *sync.Mutex
	m    map[string]map[string]int
}

func (ii *InvertedIndex) insert(word string, fileName string) {
	ii.m[word][fileName] += 1
}

func newInvertedIndex() *InvertedIndex {
	return &InvertedIndex{lock: &sync.Mutex{}, m: make(map[string]map[string]int)}
}
