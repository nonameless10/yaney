package securegraph

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

func Load() (vgs []VertexGraph, err error) {
	vgs = make([]VertexGraph, 0)

	datapath := "D:\\data\\YJS\\1\\密态知识图谱\\yaney\\yaney\\tmp\\myrdf"
	file, err := os.OpenFile(datapath, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}
	var size = stat.Size()
	fmt.Println("file size=", size)

	buf := bufio.NewReader(file)
	var count int
	count = 0
	for {
		var vg VertexGraph
		var vertexName string
		var props map[string]string
		props = make(map[string]string, 0)

		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF {
				fmt.Println("File read ok!")
				break
			} else {
				fmt.Println("Read file error!", err)
				return nil, err
			}
		}

		// 是否为uid & type
		if strings.HasPrefix(line, "uid") == true {
			//println("haha")
			line, err = buf.ReadString('\n')
			line = strings.TrimSpace(line)
			vertexName = strings.TrimSpace(strings.Split(line, "==")[0])

			//fmt.Println("vertexName:" + vertexName)
		} else {
			return nil, fmt.Errorf("error in match uid")
		}

		// 是否为properties
		line, err = buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "properties") == true {
			//fmt.Println("-----------------------properties")
			for {
				line, err = buf.ReadString('\n')
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "outEdges") == true {
					break
				}
				propName := strings.TrimSpace(strings.Split(line, "==")[0])
				propValue := strings.TrimSpace(strings.Split(line, "==")[1])
				props[propName] = propValue
				//fmt.Println("propName:" + propName + ";  propValue:" + propValue)
			}
		} else {
			return nil, fmt.Errorf("error in match properties")
		}

		vg.vertex = Vertex{
			name:  vertexName,
			props: props,
		}

		// 读取outEdge
		if strings.HasPrefix(line, "outEdges") == true {
			//fmt.Println("---------------------------outEdge")
			for {
				line, err = buf.ReadString('\n')
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "inEdges") == true {
					break
				}
				edgeType := strings.TrimSpace(strings.Split(line, "==")[0])
				peerVetexName := strings.TrimSpace(strings.Split(line, "==")[1])

				vg.outEdges = append(vg.outEdges, Edge{
					selfVertexName: vertexName,
					inout:          1,
					relationName:   edgeType,
					rank:           0,
					peerVertexName: peerVetexName,
					props:          map[string]string{},
				})

				//fmt.Println("edgeType:" + edgeType + "  ;peerVetexName:" + peerVetexName)
			}
		} else {
			return nil, fmt.Errorf("error in match outEdges")
		}

		// 读取inEdges
		if strings.HasPrefix(line, "inEdges") == true {
			//fmt.Println("---------------------------inEdges")
			for {
				line, err = buf.ReadString('\n')
				line = strings.TrimSpace(line)
				if line == "" {
					break
				}
				edgeType := strings.TrimSpace(strings.Split(line, "==")[0])
				peerVetexName := strings.TrimSpace(strings.Split(line, "==")[1])

				vg.inEdges = append(vg.inEdges, Edge{
					selfVertexName: vertexName,
					inout:          -1,
					relationName:   edgeType,
					rank:           0,
					peerVertexName: peerVetexName,
					props:          map[string]string{},
				})

				//fmt.Println("edgeType:" + edgeType + "  ;peerVetexName:" + peerVetexName)
			}
		} else {
			return nil, fmt.Errorf("error in match outEdges")
		}
		vgs = append(vgs, vg)
		count++
		//fmt.Println(count)
		if count == 100000 {
			break
		}
	}
	return
}

func CreateQuery() {
	datapath := "yaney/tmp/myrdf"
	fileRead, err := os.OpenFile(datapath, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		return
	}
	defer fileRead.Close()

	queryDataPath := "yaney/tmp/queryData"
	fileWrite, err := os.Create(queryDataPath)
	if err != nil {
		fmt.Println("Open file error!", err)
		return
	}
	defer fileWrite.Close()

	stat, err := fileRead.Stat()
	if err != nil {
		panic(err)
	}
	var size = stat.Size()
	fmt.Println("Readfile size=", size)

	w := bufio.NewWriter(fileWrite)
	buf := bufio.NewReader(fileRead)

	var count int
	count = 0
	for {
		var vertexName string

		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF {
				fmt.Println("File read ok!")
				break
			} else {
				fmt.Println("Read file error!", err)
				return
			}
		}
		// 是否为uid & type
		if strings.HasPrefix(line, "uid") == true {
			line, err = buf.ReadString('\n')
			line = strings.TrimSpace(line)
			vertexName = strings.TrimSpace(strings.Split(line, "==")[0])
			if _, err := w.WriteString(vertexName + "\n"); err != nil {
				fmt.Println(err)
			}
			count++
		}

		if count == 10000 {
			break
		}
	}
	w.Flush()
}

func RandQuery(times int) int64 {
	var querySet []string
	querySet = make([]string, 0)
	// io
	datapath := "yaney/tmp/queryData"
	file, err := os.OpenFile(datapath, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		return 0
	}
	defer file.Close()
	buf := bufio.NewReader(file)
	var count int
	count = 0
	for {
		vn, _ := buf.ReadString('\n')
		vn = strings.TrimSpace(vn)
		querySet = append(querySet, vn)
		count++
		if count == 10000 {
			break
		}
	}

	// quert test
	startTime1 := time.Now()
	masterSecret := GetSecret("master", "secret", 4096)
	keySecret := GetSecret("key", "secret", 4096)
	valSecret := GetSecret("val", "secret", 4096)
	sg := NewSecureGraph("graphdir", masterSecret, keySecret, valSecret)

	testTimes := times
	rand.Seed(time.Now().Unix())
	for {
		sg.QueryGraph(querySet[rand.Intn(10000)])
		testTimes--
		if testTimes == 0 {
			break
		}
	}
	elapsedTime1 := time.Since(startTime1) / time.Millisecond // duration in ms
	fmt.Printf("Segment finished in %d ms\n", elapsedTime1)
	sg.kv.Close()

	return int64(elapsedTime1)
}

func TestLoad() int64 {
	masterSecret := GetSecret("master", "secret", 4096)
	keySecret := GetSecret("key", "secret", 4096)
	valSecret := GetSecret("val", "secret", 4096)

	sg := NewSecureGraph("graphdir", masterSecret, keySecret, valSecret)
	startTime1 := time.Now()
	vgs, err := Load()
	if err != nil {
		fmt.Println(err)
	}
	elapsedTime1 := time.Since(startTime1) / time.Millisecond // duration in ms
	fmt.Printf("Segment finished in %d ms\n", elapsedTime1)

	startTime := time.Now()
	sg.BuildGraph(vgs)
	elapsedTime := time.Since(startTime) / time.Millisecond // duration in ms
	fmt.Printf("Segment finished in %d ms\n", elapsedTime)

	return int64(elapsedTime)
}

func RandQueryRaw(times int) int64 {
	var querySet []string
	querySet = make([]string, 0)
	// io
	datapath := "yaney/tmp/queryData"
	file, err := os.OpenFile(datapath, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		return 0
	}
	defer file.Close()
	buf := bufio.NewReader(file)
	var count int
	count = 0
	for {
		vn, _ := buf.ReadString('\n')
		vn = strings.TrimSpace(vn)
		querySet = append(querySet, vn)
		count++
		if count == 10000 {
			break
		}
	}

	// quert test
	startTime1 := time.Now()
	rg := NewRawGraph("rawgraphdir")
	testTimes := times
	rand.Seed(time.Now().Unix())
	for {
		rg.QueryGraph(querySet[rand.Intn(10000)])
		testTimes--
		if testTimes == 0 {
			break
		}
	}
	elapsedTime1 := time.Since(startTime1) / time.Millisecond // duration in ms
	rg.kv.Close()
	fmt.Printf("Segment finished in %d ms \n", elapsedTime1)

	return int64(elapsedTime1)
}

func TestLoadRaw() int64 {

	rg := NewRawGraph("rawgraphdir")
	startTime1 := time.Now()
	vgs, err := Load()
	if err != nil {
		fmt.Println(err)
		return 0
	}
	elapsedTime1 := time.Since(startTime1) / time.Millisecond // duration in ms
	fmt.Printf("Segment finished in %d ms\n", elapsedTime1)
	fmt.Printf("vgs len:%d\n", len(vgs))

	startTime := time.Now()

	rg.BuildGraph(vgs)
	elapsedTime := time.Since(startTime) / time.Millisecond // duration in ms
	fmt.Printf("Segment finished in %d ms\n", elapsedTime)

	return int64(elapsedTime)
	//vg1, _ := rg.QueryGraph("0x1364a")
	//fmt.Println(vg1)

}
