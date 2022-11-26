package main

import (
	"course_work_parallel_computing/indexer/blendedIndexer"
)

var (
	path    string
	threads int
)

func main() {
	//if flag.StringVar(&path, "path", "", "path to data that will be indexed"); path[len(path)-1:] != "/" {
	//	path += "/"
	//}

	path = "./test/"
	threads = 2

	index, err := blendedIndexer.BuildIndex(path, threads)
	if err != nil {
		return
	}

	index.PrintIndexTable()

	index.PrintSingleWordEntering("string")
}
