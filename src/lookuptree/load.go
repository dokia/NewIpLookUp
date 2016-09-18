/*
 * 加载ip库文件，排序、处理重叠ip段
 */
package lookuptree

import (
	"container/list"
	"encoding/csv"
	"os"
	"sort"
)

func Load(files []string) (l *list.List, err error) {
	// 创建切片以存储所有ip段信息
	geoindexall := make([]*GeoIndexInfo, 0, 0)
	// 循环将每个ip库文件加载进切片
	for _, file := range files {
		Infof("Load ip file: %s\n", file)
		f, err := os.Open(file)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		reader := csv.NewReader(f)
		totallines, err := reader.ReadAll()
		geoindextmp := make([]*GeoIndexInfo, len(totallines), len(totallines))
		for i, line := range totallines {
			//line := totallines[i]
			// 将ip段构建为地理位置对象
			var geoindex *GeoIndexInfo
			if len(line) >= 4 {
				geoindex = &GeoIndexInfo{IPstart: line[0],
					IPend:    line[1],
					Country:  "China",
					Province: line[2],
					City:     line[3]}
			} else if len(line) == 3 {
				geoindex = &GeoIndexInfo{IPstart: line[0],
					IPend:   line[1],
					Country: line[2]}
			} else {
				continue
			}
			geoindextmp[i] = geoindex
		}
		geoindexall = append(geoindexall, geoindextmp...)
	}

	// 按起始ip排序
	Debugln("Sorting the original data")
	sort.Sort(ByStartIp(geoindexall))
	// 处理重叠ip段，并存入list
	Debugln("Handling overlap data")
	countline := 0
	totalnum := len(geoindexall)
	l = list.New()
	for {
		if countline >= totalnum {
			break
		}
		geoindex := geoindexall[countline]
		countline++
		if geoindex == nil {
			continue
		}

		ipstartvalue, err := IpToLong(geoindex.IPstart)
		if err != nil {
			continue
		}
		ipendvalue, err := IpToLong(geoindex.IPend)
		if err != nil {
			continue
		}

		// 从后往前遍历list并在合适位置插入当前ip段
		e := l.Back()
		for {
			if e == nil {
				l.PushBack(geoindex)
				break
			}
			curgeoindex := e.Value.(*GeoIndexInfo)
			curipstartvalue, err := IpToLong(curgeoindex.IPstart)
			if err != nil {
				continue
			}
			curipendvalue, err := IpToLong(curgeoindex.IPend)
			if err != nil {
				continue
			}
			// 当前节点和新节点比较
			if curipendvalue < ipstartvalue {
				// 若新节点在当前节点之后，将新节点插入当前节点之前
				l.InsertAfter(geoindex, e)
				break
			} else if curipstartvalue <= ipstartvalue && curipendvalue >= ipendvalue {
				// 若当前节点包含新节点
				if curipstartvalue == ipstartvalue {
					if curipendvalue == ipendvalue {
						// 对于完全重合的ip段，省市覆盖国家
						if geoindex.Province != "" && geoindex.City != "" {
							e.Value.(*GeoIndexInfo).IPstart = geoindex.IPstart
							e.Value.(*GeoIndexInfo).IPend = geoindex.IPend
							e.Value.(*GeoIndexInfo).Country = geoindex.Country
							e.Value.(*GeoIndexInfo).Province = geoindex.Province
							e.Value.(*GeoIndexInfo).City = geoindex.City
						}
						break
					}
					e.Value.(*GeoIndexInfo).IPstart = LongToIp(ipendvalue + 1)
					l.InsertBefore(geoindex, e)
				} else if curipendvalue == ipendvalue {
					e.Value.(*GeoIndexInfo).IPend = LongToIp(ipstartvalue - 1)
					l.InsertAfter(geoindex, e)
				} else {
					e.Value.(*GeoIndexInfo).IPstart = LongToIp(ipendvalue + 1)
					l.InsertBefore(geoindex, e)
					e = e.Prev()
					l.InsertBefore(&GeoIndexInfo{
						IPstart:  LongToIp(curipstartvalue),
						IPend:    LongToIp(ipstartvalue - 1),
						Country:  curgeoindex.Country,
						Province: curgeoindex.Province,
						City:     curgeoindex.City}, e)
				}
				break
			} else if curipstartvalue >= ipstartvalue && curipendvalue <= ipendvalue {
				// 若新节点包含当前节点
				if curipstartvalue == ipstartvalue {
					geoindex.IPstart = LongToIp(curipendvalue + 1)
					l.InsertAfter(geoindex, e)
					break
				} else if curipendvalue == ipendvalue {
					geoindex.IPend = LongToIp(curipstartvalue - 1)
					e = e.Prev()
				} else {
					l.InsertAfter(&GeoIndexInfo{
						IPstart:  LongToIp(curipendvalue + 1),
						IPend:    geoindex.IPend,
						Country:  geoindex.Country,
						Province: curgeoindex.Province,
						City:     geoindex.City}, e)

					geoindex.IPend = LongToIp(curipstartvalue - 1)
					e = e.Prev()
					continue
				}
			}
			// 若遍历到最前端
			if e == l.Front() {
				l.PushFront(geoindex)
				break
			}
			e = e.Prev()
		}
	}
	return l, nil
}
