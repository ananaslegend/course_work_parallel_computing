package main

import (
	"bufio"
	"course_work_parallel_computing/indexer/blendedIndexer"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("\nThe program was developed as part of the \"Parallel Computing\" course \n" +
		"by Daniil Khutoryanskyi, a student of the DA-92 group\n")

	for {
		var threads int
		var path string
		var result = make(map[int]string)
		var err error

		fmt.Printf("\nDirectory with files needed to index.\n" +
			"Folder should be in the same directory with *.exe file.\n" +
			"Files should be in txt format.\n" +
			"Skip it and ./data/ will be default. \n: ")
		path, _ = reader.ReadString('\n')
		if path = validateInputString(path); path == "" {
			path = "./data/"
		}

		isCurrentDir := true
		for isCurrentDir {
			func() {
				for {
					fmt.Printf("\nThread count: ")
					threadsStr, _ := reader.ReadString('\n')
					threads, err = strconv.Atoi(validateInputString(threadsStr))
					if err != nil || threads <= 0 {
						printFailInput()
						continue
					} else {
						break
					}
				}

				start := time.Now()
				index, err := blendedIndexer.BuildIndex(path, threads)
				if err != nil {
					fmt.Printf("\nFAIL - %s\n", err)
					isCurrentDir = false
					return
				}
				duration := time.Since(start)
				result[threads] = duration.String()
				fmt.Printf("\nIndex success built in %s!\n", duration.String())

				for {
					fmt.Printf("\n1. Rebuild index" +
						"\n2. Search word " +
						"\n3. Print index table" +
						"\n4. Time result" +
						"\n5. Benchmark" +
						"\n:")

					cStr, _ := reader.ReadString('\n')
					c, err := strconv.Atoi(validateInputString(cStr))

					if err != nil {
						printFailInput()
						continue
					} else {
						switch c {
						case 1:
							for {
								fmt.Printf("\n1. Change threads count" +
									"\n2. Change directory (result table will be terminate)\n:")

								iStr, _ := reader.ReadString('\n')
								if i, _ := strconv.Atoi(validateInputString(iStr)); i == 1 {
									return
								} else if i == 2 {
									isCurrentDir = false
									return
								} else {
									printFailInput()
									continue
								}
							}
						case 2:
							for {
								fmt.Printf("\nWord: ")
								word, _ := reader.ReadString('\n')
								if word = validateInputString(word); word == "" {
									continue
								} else {
									index.PrintSingleWordEntering(word)
									break
								}
							}
						case 3:
							index.PrintIndexTable()
						case 4:
							printResult(result)
						case 5:
							var iterCount int = 1
							for {
								fmt.Printf("\nIterations\n:")
								iStr, _ := reader.ReadString('\n')
								iterCount, err = strconv.Atoi(validateInputString(iStr))
								if err != nil {
									printFailInput()
									continue
								} else {
									break
								}
							}
							benchmark(iterCount, path)
						default:
							printFailInput()
						}
					}
				}
			}()
		}
		continue
	}
}

func printResult(result map[int]string) {
	keys := make([]int, 0, len(result))
	for k := range result {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "\nThreads\t Duration")

	for _, v := range keys {
		fmt.Fprintln(w, v, "\t", result[v])
	}

	if err := w.Flush(); err != nil {
		log.Println(err)
		return
	}
}

func validateInputString(str string) string {
	return strings.Trim(str, "\n\r")
}

func printFailInput() {
	fmt.Printf("\nFAIL - Incorrect value\n" +
		"Please try again\n")
}

func benchmark(iterations int, path string) {
	threads := []int{1, 2, 3, 4, 6, 8, 10, 16, 50, 100, 500, 1000, 3000, 6000, 12500}
	result := make(map[int]time.Duration)

	for i := 1; i <= iterations; i++ {
		for _, th := range threads {
			start := time.Now()
			_, err := blendedIndexer.BuildIndex(path, th)
			if err != nil {
				fmt.Printf("\nFAIL - %s\n", err)
				return
			}
			duration := time.Since(start)

			result[th] += duration
		}
	}
	for index, durationSum := range result {
		result[index] = durationSum / time.Duration(iterations) // average value
	}

	resultStr := func() map[int]string {
		var resultStr = make(map[int]string, len(result))
		for i, v := range result {
			resultStr[i] = v.String()
		}
		return resultStr
	}

	printResult(resultStr())
}
