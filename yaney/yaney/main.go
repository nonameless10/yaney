package main

import (
	"fmt"
	"os/exec"
	"runtime"
)

// 明文构建
func testLoad() {
	var i int64
	var res []int64
	var times int64
	times = 1
	// 只循环一次
	for i = 0; i < times; i++ {
		rtn := securegraph.TestLoad()
		res = append(res, rtn)
		fmt.Printf("------done %d %d\n", i, times)
		if i != times-1 {
			removeDir("yaney/graphdir")
		}
	}

	var sum int64
	for _, v := range res {
		sum += v
	}
	fmt.Println(res)
	fmt.Printf("avg res: %d", sum/times)

}

func testQueryOne() {
	masterSecret := securegraph.GetSecret("master", "secret", 4096)
	keySecret := securegraph.GetSecret("key", "secret", 4096)
	valSecret := securegraph.GetSecret("val", "secret", 4096)

	sg := securegraph.NewSecureGraph("graphdir", masterSecret, keySecret, valSecret)
	//fmt.Println("hh")
	vg1, _ := sg.QueryGraph("0x7")
	fmt.Println(vg1)
}

func testQueryManyTimes() {
	var i int64
	var res []int64
	var times int64

	times = 50
	for i = 0; i < times; i++ {
		rtn := securegraph.RandQuery(10000)
		//fmt.Println("-------done\n")
		res = append(res, rtn)
		fmt.Printf("------done %d\n", i)
		//time.Sleep(4 * time.Duration(time.Second))
	}
	var sum int64
	for _, v := range res {
		sum += v
	}
	fmt.Println(res)
	fmt.Printf("avg res: %d\n", sum/times)
}

func testRawLoad() {
	var i int64
	var res []int64
	var times int64

	times = 1
	for i = 0; i < times; i++ {
		rtn := securegraph.TestLoadRaw()
		res = append(res, rtn)
		fmt.Printf("------done %d\n", i)
		if i != times-1 {
			removeDir("/home/l1nkkk/project/mime/yaney/rawgraphdir")
		}
	}
	var sum int64
	for _, v := range res {
		sum += v
	}
	fmt.Println(res)
	fmt.Printf("avg res: %d", sum/times)
}

func testRawQueryOne() {
	rg := securegraph.NewRawGraph("graphdir")
	vg1, _ := rg.QueryGraph("0x13655")
	fmt.Println(vg1)
}

func testRawQueryManyTimes() {
	var i int64
	var res []int64
	var times int64

	times = 50
	for i = 0; i < times; i++ {
		rtn := securegraph.RandQueryRaw(10000)
		res = append(res, rtn)
		fmt.Printf("------done %d\n", i)

	}
	var sum int64
	for _, v := range res {
		sum += v
	}
	fmt.Println(res)
	fmt.Printf("avg res: %d", sum/times)
}

func removeDir(dirpath string) {
	fmt.Println("Remove " + dirpath)
	cmd := exec.Command("/bin/bash", "-c", "rm -rf "+dirpath)
	cmd.Output()
}

func main() {
	runtime.GOMAXPROCS(1)
	//testQueryManyTimes()
	//testRawLoad()
	testLoad()
	//fmt.Printf("test\n")
	testQueryOne()
	//testRawQueryOne()
	//testRawQueryManyTimes()

	//testQueryManyTimes()

}
