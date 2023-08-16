package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// 肉品架構
type Meat struct {
	Name string        `json:"name" gorm:"column:name"`
	Num  int           `json:"num" gorm:"column:num"`
	Time time.Duration `json:"time" gorm:"column:time"`
}

// 肉品庫存資訊
var stock []Meat = []Meat{
	{Name: "牛肉", Num: 10, Time: 1},
	{Name: "豬肉", Num: 7, Time: 2},
	{Name: "雞肉", Num: 5, Time: 3},
}

func main() {
	var lock sync.Mutex   // 宣告Lock 用以資源佔有與解鎖
	var wg sync.WaitGroup // 宣告WaitGroup 用以等待執行序
	go work(&wg, &lock, "E")
	go work(&wg, &lock, "D")
	go work(&wg, &lock, "C")
	go work(&wg, &lock, "B")
	go work(&wg, &lock, "A")
	wg.Add(5) // 記數器（五個員工
	wg.Wait() // 等待計數器歸零
}

// 印出資訊
func log(worker string, status bool, name string) {
	sta := "取得"
	if status {
		sta = "處理完"
	}
	fmt.Println("┌─────────────────────────────────────────────────────────┐")
	fmt.Println("│", worker, "在", time.Now().Format("2006-01-02 15:04:05"), sta, name)
	if len(stock) > 0 {
		fmt.Println("│────────────────────────── 庫存 ─────────────────────────│")
		for _, meat := range stock {
			fmt.Println("│", meat.Name, ":", meat.Num)
		}
	}
	fmt.Println("└─────────────────────────────────────────────────────────┘")
}

func work(wg *sync.WaitGroup, lock *sync.Mutex, worker string) {
	defer wg.Done() //wg 計數器-1
	for len(stock) > 0 {
		lock.Lock()
		i := rand.Intn(len(stock)) // 隨機取肉品
		curr := stock[i]
		stock[i].Num--
		if stock[i].Num == 0 {
			stock = append(stock[:i], stock[i+1:]...)
		}
		log(worker, false, curr.Name)
		lock.Unlock()
		time.Sleep(curr.Time * time.Second)
		log(worker, true, curr.Name)
	}
}
