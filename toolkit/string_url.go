package toolkit

import (
	"strings"
)

func FindIpAddress(url string) (int, int) {
	// http://113.200.102.123:12401
		// get ip
	idx1 := strings.Index(url, `//`)
	if idx1 > 0 {
		idx1 += 2
		idx2 := strings.Index(url[idx1:], `:`) + idx1
		if idx2 > idx1 {
			// we got ip address
//			fmt.Println("ip=%v", url[idx1: idx2])

			return idx1, idx2
		} else {
//			fmt.Println("idx1=%v, idx2 =%v", idx1, idx2)
		}
	} else {
//		fmt.Println("idx1 < 0")
	}		
	return -1,-1
}

func GetProtocolIpPort(url string) (string, string) {
	// http://113.200.102.123:12401
	// https://113.200.102.123:12401
	idx1 := strings.Index(url, `//`)
	if idx1 > 0 {
		idx1 += 2
		idx2 := strings.Index(url[idx1:], `/`) + idx1
		if idx2 > idx1 {
			return url[0:idx2], url[idx2:]
		} 
	}
	return "", ""
}