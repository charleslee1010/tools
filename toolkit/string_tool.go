package toolkit

import (
	"strings"
	"fmt"
//	"errors"
	"strconv"
)

// convert kv pair string into map
// k1:v1,k2:v2,k3,k4:v4....
func GetKVMap(f string) map[string]interface{} {
	r := make(map[string]interface{})
	
	t := strings.Split(f, ",")
	if t == nil && len(t) == 0 {
		return r
	}
	
	// find the first :
	for _, v := range t {
		if v == "" {
//			fmt.Println("empty kv")
			continue
		}
		
		idx := strings.Index(v, ":")
		if idx < 0 {
			// not present, v is treated as key and value is ""
			fmt.Println("no colon")
			r[v] = ""
		} else if idx == 0 {
			// first char is :, v is treated as invalid, ignore
			fmt.Println("invalid kv, first colon")
		} else 	{	
			r[v[0:idx]] = v[idx+1:]
		}
	}
	return r
}

func SplitAndTrim(v string, sep string) []string {
	t := strings.Split(v, sep)
	if t == nil {
		return nil
	}
	nt := make([]string, 0)
	for _, nv := range t {
		tmp := strings.TrimSpace(nv)
		if tmp != "" {
			nt = append(nt, tmp)
		}
	}
	return nt
}


func SplitAndTrimSpace(v string, sep string) []string {
	
	if sep == "" || v == "" {
		return []string {v}
	}
	
	n := strings.Count(v, sep)
	lensep := len(sep)
	
	s := make([]string, n+1)
	left := v
	for i := 0; i < n+1; i ++ {
		idx := strings.Index(left, sep)
		if idx == 0 {
			s[i] = ""
		} else if idx > 0 {
			s[i] = strings.TrimSpace(left[0:idx])
		} else {
			s[i] = strings.TrimSpace(left)		
		}
		
		// calculate what's left after lensep and token are removed
		if idx + lensep == len(left) {
			left = ""
		} else {			
			left = left[idx+lensep:]
		}		
	}
	
	return s
}


func String2IntSlice(v string, sep string) ([]int, error) {

	t := strings.Split(v, sep)
	
	intSlice := make([]int, 0)
	
	for i := 0; i < len(t); i++ {
		if iv, err := strconv.Atoi(t[i]); err == nil {
			intSlice = append(intSlice, iv)  
		} else {
			return nil, err
		}
	}
	
	return intSlice, nil
}

