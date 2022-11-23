package main

import (
	"course_work_parallel_computing/indexer/blendedIndexer"
)

func main() {
	//threadCount := 34

	index, err := blendedIndexer.BuildIndex("./test/", 2)
	if err != nil {
		return
	}

	index.Print()
}
