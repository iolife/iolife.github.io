package main

import (
	"log"
	"net/http"
	"sync"

	_ "net/http/pprof"
	"time"
)

func Chunk[T any](collection []T, size int) [][]T {
	if size <= 0 {
		panic("Second parameter must be greater than 0")
	}

	chunksNum := len(collection) / size
	if len(collection)%size != 0 {
		chunksNum += 1
	}

	result := make([][]T, 0, chunksNum)

	for i := 0; i < chunksNum; i++ {
		last := (i + 1) * size
		if last > len(collection) {
			last = len(collection)
		}
		result = append(result, collection[i*size:last])
	}

	return result
}
func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	var depotID int64 = 1
	GetChildIdArray := func(ids []int64) []int64 {
		for i := 0; i < 1000; i++ {
			ids = append(ids, int64(i))
		}
		return ids
	}
	getPermissionDataByNodeIds := func(nodeIds []int64) {
		time.Sleep(1 * time.Millisecond)
	}

	chunkIDs := [][]int64{{depotID}}
	ch := make(chan struct{}, 100)

	for len(chunkIDs) != 0 {
		var mu sync.Mutex
		var tmpIds [][]int64
		var wg sync.WaitGroup
		wg.Add(len(chunkIDs) * 2)
		for _, v := range chunkIDs {
			time.Sleep(1 * time.Second)
			chunkID := v
			ch <- struct{}{}
			go func() {
				defer wg.Done()
				getPermissionDataByNodeIds(chunkID)
				<-ch
			}()

			ch <- struct{}{}
			go func() {
				defer wg.Done()
				ids := GetChildIdArray(chunkID)
				mu.Lock()
				tmpIds = append(tmpIds, Chunk(ids, 100)...)
				mu.Unlock()
				<-ch
			}()
		}
		wg.Wait()
		chunkIDs = tmpIds
	}

}

