package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type LogProcess struct {
	rc    chan []byte
	wc    chan string
	read  Reader
	write Writer
}

type Reader interface {
	Read(rc chan []byte)
}

type Writer interface {
	Write(wc chan string)
}

type ReadFromFile struct {
	path string //读取文件的路径
}

func (r *ReadFromFile) Read(rc chan []byte) {
	//读取模块
	//打开文件
	f, err := os.Open(r.path)
	if nil != err {
		panic(fmt.Sprintf("open file error: %s", err.Error()))
	}

	//从文件末尾开始逐行读取文件内容
	f.Seek(0, 2)
	rd := bufio.NewReader(f)

	for {
		line, err := rd.ReadBytes('\n')
		fmt.Println(line)
		if err == io.EOF {
			time.Sleep(500 * time.Millisecond)
		} else if nil != err {
			panic(fmt.Sprintf("ReadBytes error:%s", err.Error()))
		}

		rc <- line
	}
}

func (l *LogProcess) Process() {
	//解析模块
	data := <-l.rc
	l.wc <- strings.ToUpper(string(data))
}

type WriteToInfluxDB struct {
	influxDBsn string //influxdb source
}

func (w *WriteToInfluxDB) Write(wc chan string) {
	//写入模块
	fmt.Print(<-wc)
}

func main() {
	r := &ReadFromFile{
		path: "./access.log",
	}

	w := &WriteToInfluxDB{
		influxDBsn: "username&passwd",
	}

	lp := &LogProcess{
		rc: make(chan []byte),
		wc: make(chan string),
		read: r,
		write: w,
	}


	go lp.read.Read(lp.rc)
	go lp.Process()
	go lp.write.Write(lp.wc)

	time.Sleep(1 * time.Second)
}
