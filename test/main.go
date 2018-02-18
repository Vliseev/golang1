package main

import (
	"fmt"
	"os"
	"bufio"
	"io"
)

const filePath string = "/home/vad/GO/src/coursera/golang1/hw3_bench/data/users.txt"

func test() {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	for{
		line,err:=reader.ReadString('\n')
		if err!=nil {
			if err==io.EOF{
				break
			}else{
				return
			}
		}

		fmt.Println(line)
	}
}

func main()  {
	//fmt.Println(strings.Replace("AnthonyOlson@Zooxo.org", "@", "[at]", -1))
	test()
}
