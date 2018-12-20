package main

import (
	"strings"
	"time"
)

type LogProcess struct {
	rc    chan string
	wc    chan string
	read  Reader
	write Writer
}

type Reader interface {
	Read(rc chan string)
}

type Writer interface {
	Write(wc chan string)
}

type ReadFromFile struct {
	path string //读取文件的路径
}

func (r *ReadFromFile) Read(rc chan string) {
	//读取模块
	line := "message"
	rc <- line
}

func (l *LogProcess) Process() {
	//解析模块
	data := <-l.rc
	l.wc <- strings.ToUpper(data)
}

type WriteToInfluxDB struct {
	influxDBsn string //influxdb source
}

func (w *WriteToInfluxDB) Write(wc chan string) {
	//写入模块
	println(<-wc)
}

func main() {
	r := &ReadFromFile{
		path: "/tmp/access.log",
	}

	w := &WriteToInfluxDB{
		influxDBsn: "username&passwd",
	}

	lp := &LogProcess{
		rc: make(chan string),
		wc: make(chan string),
		read: r,
		write: w,
	}


	go lp.read.Read(lp.rc)
	go lp.Process()
	go lp.write.Write(lp.wc)

	time.Sleep(1 * time.Second)
}