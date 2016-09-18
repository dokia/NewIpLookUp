/*
 * 排序GeoIndexInfo对象
 */
package lookuptree

type ByStartIp []*GeoIndexInfo

func (b ByStartIp) Len() int {
	return len(b)
}

func (b ByStartIp) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b ByStartIp) Less(i, j int) bool {
	v1, _ := IpToLong(b[i].IPstart)
	v2, _ := IpToLong(b[j].IPstart)
	return v1 < v2
}
