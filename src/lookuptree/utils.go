package lookuptree

import (
	"errors"
	"strconv"
	"strings"
)

/*
 * 获得IP地址第level层的整型数(从前往后数，0至3位)
 */
func GetIpSection(ip string, level int) (ipsec int, err error) {
	ipsecs := strings.Split(ip, ".")

	if level < 0 || level >= len(ipsecs) {
		err = errors.New("Wrong index when parsing ip.")
		return
	}
	return strconv.Atoi(ipsecs[level])
}

/*
 * 将ip地址转换为64位整型
 */
func IpToLong(ip string) (result int64, err error) {
	for i := 0; i < 4; i++ {
		tmp, err := GetIpSection(ip, i)
		if err != nil {
			return int64(0), err
		}
		result += int64(tmp << uint8((3-i)*8))
	}
	return result, nil
}

/*
 * 将64位整型转换为ip地址
 */
func LongToIp(ip int64) (result string) {
	for i := 0; i < 4; i++ {
		tmp := ip & 0xff
		ip = ip >> 8
		if i == 0 {
			result = strconv.Itoa(int(tmp))
		} else {
			result = strconv.Itoa(int(tmp)) + "." + result
		}
	}
	return result
}

/*
 * 将level以下的部分都修改为255
 */
func EnlargeIP(ip string, level int) (result string) {
	ipsecs := strings.Split(ip, ".")
	result = ipsecs[0]
	for i := 1; i < 4; i++ {
		if i > level {
			result += ".255"
		} else {
			result += "." + ipsecs[i]
		}
	}
	return
}
