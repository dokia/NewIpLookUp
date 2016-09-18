/*
 * IP查找库
 * 读入本地.csv文件，每行格式为:
 *
 *			<起始ip,终止ip,省,市>
 * 		或
 *			<起始ip,终止ip,国家/地区>
 *
 * 并通过http接口提供查询服务，查询方式为
 *			http://host:port/request_ip
 */

package main

import (
	"flag"
	"fmt"
	"log"
	"lookuptree"
	"net/http"
	"strconv"
	"strings"
)

// ip查找树
var tree *lookuptree.LookUpTree

// ip库本地文件组，.csv文件，以";"隔开
var ipfile *string = flag.String("f", "ip.csv;foreign.csv", "IP file in csv format, separated by\";\"")

// 服务器监听端口
var port *int = flag.Int("p", 8080, "Listen port for local server")

func Init() {
	// 解析传入参数
	flag.Parse()
	lookuptree.Infof("Server listen port: %d\n", *port)
	// 创建查找树
	tree = lookuptree.NewLookUpTree()
	// 解析ip库文件组
	files := strings.Split(*ipfile, ";")
	// 加载文件组，排序并处理重叠的ip段
	l, err := lookuptree.Load(files)
	if err != nil {
		log.Fatal(err)
	} else {
		// 插入ip查找树
		for iter := l.Front(); iter != nil; iter = iter.Next() {
			tree.Insert(iter.Value.(*lookuptree.GeoIndexInfo))
		}
	}
	lookuptree.Infoln("Load file finished, ready to serve...")
}

func handler(w http.ResponseWriter, r *http.Request) {
	ip := r.URL.Path[1:]
	lookuptree.Debugf("Look up ip: %s\n", ip)
	// 查找请求ip
	result, err := tree.Search(ip)
	if err != nil {
		fmt.Fprintf(w, "Error: Cannot find matched info")
	} else {
		fmt.Fprintf(w, "Country/Region: "+result.Country+"\nProvince: "+result.Province+"\nCity: "+result.City)
	}
}

func main() {
	Init()
	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+strconv.Itoa(*port), nil)
}
