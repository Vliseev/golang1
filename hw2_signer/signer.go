package main

import (
	"sync"
	"strconv"
	"bytes"
	"sort"
)

type pair struct {
	first  string
	second string
}


func parMultiHash(dataString string, out chan interface{}, waitJob *sync.WaitGroup){
	defer waitJob.Done()

	const numWorker = 6

	wg := &sync.WaitGroup{}
	arrayData := make([]pair, numWorker)

	for i := 0; i < numWorker; i++ {
		wg.Add(1)
		arrayData[i].first = strconv.Itoa(i) + dataString
		go func(data_pair []pair, i int, wg *sync.WaitGroup) {
			defer wg.Done()
			data_pair[i].second = DataSignerCrc32(data_pair[i].first)
		}(arrayData, i, wg)
	}
	wg.Wait()

	var buffer bytes.Buffer

	for idx := range arrayData {
		buffer.WriteString(arrayData[idx].second)
	}

	out <- buffer.String()
}

func MultiHash(in, out chan interface{}) {

	waitJob := &sync.WaitGroup{}

	for data := range in {
		waitJob.Add(1)
		go 	parMultiHash(data.(string),out,waitJob)
	}
	waitJob.Wait()

}


func parSingleHash(dataString string,md5Result string, out chan interface{}, waitJob *sync.WaitGroup){
	defer waitJob.Done()

	const numWorker = 2

	arrayData := make([]pair, numWorker)
	result := make(chan struct{})

	arrayData[0].first=dataString
	arrayData[1].first=md5Result

	go func(data_pair []pair) {
		data_pair[1].second = DataSignerCrc32(data_pair[1].first)
		result <- struct{}{}
	}(arrayData)

	arrayData[0].second = DataSignerCrc32(arrayData[0].first)

	<-result

	out <- arrayData[1].second + "~" + arrayData[0].second
}


func SingleHash(in, out chan interface{}) {


	waitJob := &sync.WaitGroup{}

	for data := range in {
		waitJob.Add(1)
		dataString:=strconv.Itoa(data.(int))
		md5Result := DataSignerMd5(dataString)

		go 	parSingleHash(md5Result,dataString,out,waitJob)
	}
	waitJob.Wait()

}

func CombineResults(in, out chan interface{}) {
	var accmt []string
	for val := range in {
		accmt = append(accmt, val.(string))
	}

	sort.Strings(accmt)

	var buffer bytes.Buffer

	for i :=0;i<len(accmt)-1;i++{
		buffer.WriteString(accmt[i])
		buffer.WriteString("_")
	}
	buffer.WriteString(accmt[len(accmt)-1])

	out<-buffer.String()

}

func ExecutePipeline(jobs []job) {
	var numWorker = len(jobs)

	var chans []chan interface{}

	for i := 0; i < numWorker+1; i++ {
		chans = append(chans, make(chan interface{}))
	}

	for i := 0; i < len(jobs); i++ {
		tmtJob := jobs[i]
		go job(func(in, out chan interface{}) {
			tmtJob(in, out)
			close(out)
		})(chans[i], chans[i+1])
	}

	<-chans[numWorker]

}
