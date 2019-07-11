package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
)

func main() {
	//inputData := []int{0, 1, 1, 2, 3, 5, 8}
	//inputData := []int{0,1}
	//testResult := ""
	//hashSignJobs := []job{
	//	job(func(in, out chan interface{}) {
	//		for _, fibNum := range inputData {
	//			out <- fibNum
	//		}
	//	}),
	//	job(SingleHash),
	//	job(MultiHash),
	//	job(CombineResults),
	//
	//	job(func(in, out chan interface{}) {
	//		dataRaw := <-in
	//		data, ok := dataRaw.(string)
	//		if !ok {
	//			fmt.Println("cant convert result data to string")
	//		}
	//		testResult = data
	//	}),
	//}
	//
	//start := time.Now()
	//
	//ExecutePipeline(hashSignJobs...)
	//
	//end := time.Since(start)
	//
	//fmt.Println(end)
}

const TH  = 6

func ExecutePipeline(jobs... job)  {
	wg := &sync.WaitGroup{}
	in := make(chan interface{})

	for _, jb := range jobs{
		wg.Add(1)
		out := make(chan interface{})
		go func(jobFunc job, in, out chan interface{}, wg *sync.WaitGroup) {
			defer wg.Done()
			defer close(out)
			jobFunc(in, out)
		}(jb, in, out, wg)

		in = out
	}

	defer wg.Wait()
}

func SingleHash(in, out chan interface{}){
	wg := &sync.WaitGroup{}
	insideCh := make(chan string)
	for dataIn := range in{
		wg.Add(1)
		data := fmt.Sprintf("%v", dataIn)
		md5 := DataSignerMd5(data)
		go func(data string, wg *sync.WaitGroup, insideCh chan string, md5 string) {
			defer wg.Done()

			/////////////////////////////////////////////////////////////////////
			// Парралелим СRC32
			insideInsideCh := make(chan struct{th int; hash string})
			insideInsideWg := &sync.WaitGroup{}

			// создаем структуры для горутина
			md5Struct := struct {
				th int
				hash string
			}{}
			md5Struct.th = 1
			md5Struct.hash = md5
			//
			dataStruct := struct {
				th int
				hash string
			}{}
			dataStruct.th = 0
			dataStruct.hash = data

			dataSl := []struct{th int; hash string}{dataStruct, md5Struct}

			for _, strc := range dataSl{
				insideInsideWg.Add(1)
				go func(iiCh chan struct{th int; hash string}, iiWg *sync.WaitGroup, data struct{th int; hash string}) {
					defer iiWg.Done()
					data.hash = DataSignerCrc32(data.hash)
					iiCh <- data
				}(insideInsideCh, insideInsideWg, strc)
			}

			go func(wg *sync.WaitGroup, insideInsideCh chan struct{th int; hash string}) {
				defer close(insideInsideCh)
				wg.Wait()
			}(insideInsideWg, insideInsideCh)


			for hash := range insideInsideCh{
				dataSl[hash.th].hash = hash.hash
			}
			/////////////////////////////////////////////////////////////////////
			//старый варик / надо параллелить crc32
			//crc32Md5 := DataSignerCrc32(md5)
			//crc32Data := DataSignerCrc32(data)
			//res := crc32Data + "~" + crc32Md5

			insideCh <- dataSl[0].hash + "~" + dataSl[1].hash

		}(data, wg, insideCh, md5)
	}

	go func(wg *sync.WaitGroup, insideCh chan string) {
		defer close(insideCh)
		wg.Wait()
	}(wg, insideCh)

	for res := range insideCh{
		out <- res
	}
}

func MultiHash(in, out chan interface{}){
	outsideWg := &sync.WaitGroup{}
	insideWg := &sync.WaitGroup{}
	outsideCh := make(chan string)

	for sHash := range in{
		outsideWg.Add(1)
		go func(outsideChan chan string, outsideWg *sync.WaitGroup, data string) {
			defer outsideWg.Done()

			// канал для вывода с
			insideCh  := make(chan struct{th int; subHash string})
			// дата из предыдущего канала
			insideWg.Add(TH)
			for i := 0; i < TH; i++{
				go CalculateSubHash(i, data, insideCh, insideWg)
			}
			// горутин который ждет завержения всех горутин и закрывает канал
			go func(insideCh chan struct{th int; subHash string}, insideWg *sync.WaitGroup) {
				defer close(insideCh)
				insideWg.Wait()
			}(insideCh, insideWg)
			// чтение из внутреннего канала
			// складываем хеши в слайс
			subHashes := make([]struct{th int; subHash string}, 0)
			for subHash := range insideCh{
				subHashes = append(subHashes, subHash)
			}
			// сортируем по th
			sort.Slice(subHashes, func(i, j int) bool {
				return subHashes[i].th < subHashes[j].th
			})
			//формируем строку mhHash
			res := ""
			for _, item := range subHashes{
				res += item.subHash
			}
			// посылаем в верном порядке
			outsideChan <- res

		}(outsideCh, outsideWg, sHash.(string))
	}

	// горутин который ждет завержения всех горутин и закрывает канал outsideCh
	go func(outsideCh chan string, outsideWg *sync.WaitGroup) {
		defer close(outsideCh)
		outsideWg.Wait()
	}(outsideCh, outsideWg)

	// чтение из внешнего канала
	for hash := range outsideCh{
		out <- hash
	}
}

func CalculateSubHash(th int, data string, insideCh chan struct{th int; subHash string}, insideWg *sync.WaitGroup){
	defer insideWg.Done()
	result := struct {
		th int
		subHash string
	}{}
	result.th = th
	result.subHash = DataSignerCrc32(strconv.Itoa(th) + data)
	insideCh <- result
}

func CombineResults(in, out chan interface{}){
	sl := make([]string, 0)
	for mhHash := range in{
		sl = append(sl, mhHash.(string))
	}
	sort.Strings(sl)
	res := strings.Join(sl, "_")
	out <- res
}