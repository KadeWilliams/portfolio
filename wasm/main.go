package wasm

import (
	"math/rand"
	"syscall/js"
	"time"
)

func main() {
	// js.Global().Set("goSort", js.FuncOf(sort))//<-make(chan bool)
	js.Global().Set("goSort", js.FuncOf(sort));<- make(chan bool) // keep alive
	// js.Global().Get("console").Call("log", "WASM is working!");<-make(chan bool)
}

func sort(this js.Value, args []js.Value) interface{} {
	arr := generateRandomArray(10) // 10 elements
	bubbleSort(arr, js.Global().Get("updateSortBars"))
	return nil
}

func bubbleSort(arr []int, callback js.Value) {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-i-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
			}
			// update visualization after each step

			callback.Invoke(js.ValueOf(arr))
			time.Sleep(100 * time.Millisecond) // for animation
		}
	}
}

func generateRandomArray(size int) []int {
	arr := make([]int, size)
	for i := range arr {
		arr[i] = rand.Intn(100) + 1 // 1-100
	}
	return arr 
}