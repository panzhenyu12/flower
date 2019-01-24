package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/sony/sonyflake"
)

var numchan = make(chan int64, 1000)

func test() {
	var st sonyflake.Settings
	//local := time.Local
	t, _ := time.Parse("2006-01-02", "2018-01-01")
	st.StartTime = t
	//st.MachineID = "awdaw"
	sf := sonyflake.NewSonyflake(st)

	if sf == nil {
		panic("sonyflake not created")
	}
	id, err := sf.NextID()
	if err != nil {
		fmt.Println(err)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(sonyflake.Decompose(id))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
func test1(nodeid int64) {
	// Create a new Node with a Node number of 1
	fmt.Printf("nodeid ID: %s\n", nodeid)
	node, err := snowflake.NewNode(nodeid)
	if err != nil {
		fmt.Println(err)
		return
	}
	for {

		// Generate a snowflake ID.
		id := node.Generate()
		// Print out the ID in a few different ways.
		fmt.Printf("Int64  ID: %d\n", id)
		bt, _ := id.MarshalJSON()
		fmt.Println(string(bt))
		//numchan <- id.Int64()
		fmt.Printf("String ID: %s\n", id)
		fmt.Printf("Base2  ID: %s\n", id.Base2())
		fmt.Printf("Base64 ID: %s\n", id.Base64())

		// Print out the ID's timestamp
		fmt.Printf("ID Time  : %d\n", id.Time())

		// Print out the ID's node number
		fmt.Printf("ID Node  : %d\n", id.Node())

		// Print out the ID's sequence number
		fmt.Printf("ID Step  : %d\n", id.Step())

		// Generate and print, all in one.
		fmt.Printf("ID       : %d\n", node.Generate().Int64())
	}
	// fmt.Printf("String ID: %s\n", id)
	// fmt.Printf("Base2  ID: %s\n", id.Base2())
	// fmt.Printf("Base64 ID: %s\n", id.Base64())

	// // Print out the ID's timestamp
	// fmt.Printf("ID Time  : %d\n", id.Time())

	// // Print out the ID's node number
	// fmt.Printf("ID Node  : %d\n", id.Node())

	// // Print out the ID's sequence number
	// fmt.Printf("ID Step  : %d\n", id.Step())

	// // Generate and print, all in one.
	// fmt.Printf("ID       : %d\n", node.Generate().Int64())
}

func RandInt(min, max int64) int64 {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Int63n(max-min)
}
func RandHanzi() string {
	a := make([]rune, 3)
	for i := range a {
		a[i] = rune(RandInt(19968, 40869))
	}
	return string(a)
}

func looptest() {
	for i := 1; i < 10; i++ {
		id := i
		go test1(int64(id))
	}
	go func() {
		idmap := make(map[int64]int64)
		for {
			id := <-numchan
			//fmt.Println(id)
			_, ok := idmap[id]
			if ok {
				fmt.Println("existid", id)
			} else {
				idmap[id] = id
			}
		}

	}()
}
