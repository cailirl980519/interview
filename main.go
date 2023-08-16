package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

// 肉品架構
type Meat struct {
	Name  string        `json:"name" gorm:"column:name"`
	Num   int           `json:"num" gorm:"column:num"`
	Count int           `json:"count" gorm:"column:count"`
	Time  time.Duration `json:"time" gorm:"column:time"`
}

// 肉品庫存資訊
var stock []Meat = []Meat{
	{Name: "牛肉", Num: 10, Count: 0, Time: 1},
	{Name: "豬肉", Num: 7, Count: 0, Time: 2},
	{Name: "雞肉", Num: 5, Count: 0, Time: 3},
}

var ch (chan Meat) = make(chan Meat)

func main() {
	var wg sync.WaitGroup // 宣告WaitGroup 用以等待執行序
	go work(&wg, "E")
	go work(&wg, "D")
	go work(&wg, "C")
	go work(&wg, "B")
	go work(&wg, "A")
	go putMeatToChan()
	wg.Add(5) // 記數器（五個員工
	wg.Wait() // 等待計數器歸零
}

// 印出資訊
func log(worker string, status bool, name string) {
	sta := "取得"
	if status {
		sta = "處理完"
	}
	str := []string{
		" ┌─────────────────────────────────────────────────────────┐\n",
	}
	str = append(str, "│", worker, "在", time.Now().Format("2006-01-02 15:04:05"), sta, name, "\n")
	str = append(str, "│────────────────────────── 庫存 ─────────────────────────│\n")
	for _, meat := range stock {
		str = append(str, "│", meat.Name, ":", strconv.Itoa(meat.Num))
	}
	str = append(str, "\n │──────────────────────── 生產線中 ───────────────────────│\n")
	for _, meat := range stock {
		str = append(str, "│", meat.Name, ":", strconv.Itoa(meat.Count))
	}
	str = append(str, "\n └─────────────────────────────────────────────────────────┘\n")
	fmt.Print(strings.Join(str, " "))
}

// 將肉品放置Channel(生產線)
func putMeatToChan() {
	suffle := []int{}
	for i := 0; i < 3; i++ {
		suffle = append(suffle, i)
	}
	for len(suffle) > 0 {
		i := rand.Intn(len(suffle)) // 隨機取肉品
		meat := stock[i]
		stock[i].Num--
		stock[i].Count++
		ch <- meat
		if stock[i].Num == 0 {
			suffle = append(suffle[:i], suffle[i+1:]...)
		}
	}
	close(ch)
}

func work(wg *sync.WaitGroup, worker string) {
	defer wg.Done() //wg 計數器-1
	for len(stock) > 0 {
		meat, ok := <-ch //拿取肉品
		if ok {
			// 從生產線計數移除
			for i, v := range stock {
				if v.Name == meat.Name {
					stock[i].Count--
					break
				}
			}
			log(worker, false, meat.Name)
			time.Sleep(meat.Time * time.Second)
			log(worker, true, meat.Name)
		} else {
			return
		}
	}
}
