package toolkit

import (
	"strings"
	"strconv"
	"errors"
	"fmt"
	"reflect"
	"os"
	"path/filepath"
)

func GetExecFileName() string {
	if len(os.Args) > 0 {
		return filepath.Base(os.Args[0])
	}	
	return ""
}

//func GetMapValueBool(mp map[string]interface{}, key string) bool {
//	if v, pres := mp[key]; pres {
//		if v != nil && v.(string) != "0" {
//			return true
//		}
//	}
//	return false
//}

//func GetMapValueString(mp map[string]interface{}, key string) string {
//	if v, pres := mp[key]; pres && v != nil {
//		if vs, ok := v.(string); ok {
//			return strings.TrimSpace(vs)
//		}
//		
//		return fmt.Sprint(v)
//	}
//	return ""
//}

func GetMapValueString(mp map[string]interface{}, key string, dft string) (string, error) {
	if v, pres := mp[key]; pres {
		if v == nil {
			return dft,errors.New("key:" + key + " value is nil")
		}
		if vs, ok := v.(string); ok {
			return strings.TrimSpace(vs), nil
		} 
		
		if vs, ok := v.(float64); ok {
//			return strconv.FormatFloat(vs, ), nil
			return strconv.FormatInt(int64(vs), 10), nil
		} 
		
		if vs, ok := v.(int64); ok {
			return strconv.FormatInt(vs, 10), nil
		}

		if vs, ok := v.(int); ok {
			return strconv.Itoa(vs), nil
		}

		return dft, fmt.Errorf("key:%s, type of value unknown:%s", key, reflect.TypeOf(v).String())
		
	} else {
		return dft, errors.New("key:" + key + " not found")
	}
}

func SetMapValue(mp map[string]interface{}, key1, key2 string, newv interface{}) (error) {
	
	var m1 map[string]interface{} = nil
	
	if v, pres := mp[key1]; pres && v != nil {
		if vm, ok := v.(map[string]interface{}); ok {
			m1 = vm
		} else {
			// not map, do nothing
			return errors.New("invalid key type")
		}
	} else {
		// key1 is not found
		m1 = make(map[string]interface{})
	}

	m1[key2] = newv
	mp[key1] = m1
	return nil
}


func SetMapValueString1(mp map[string]interface{}, key1, key2, v string) (map[string]interface{}, error) {
	
	var m1 map[string]interface{} = nil
	
	if v, pres := mp[key1]; pres && v != nil {
		if vm, ok := v.(map[string]interface{}); ok {
			m1 = vm
		} else {
			// not map, do nothing
			return mp, errors.New("invalid key type")
		}
	} else {
		// key1 is not found
		m1 = make(map[string]interface{})
	}

	m1[key2] = v
	mp[key1] = m1
	return mp, nil
}

func DeleteMapParameter(mp map[string]interface{}, key1, param string) (error) {
	// delete map[key1].param
	
	if v, pres := mp[key1]; pres && v != nil {
		if vm, ok := v.(map[string]interface{}); ok && vm != nil{
			// try to find param
			delete(vm, param) 
			return nil
		} else {
			// not map, do nothing
			return errors.New("not a map")
		}
	} else {
		// key1 is not found || v is nil
		return errors.New("key1 not found or not a map")
	}
}


func DeleteMapKey(mp map[string]interface{}, key1, param string) (error) {
	// delete map[key1].param
	
	if v, pres := mp[key1]; pres && v != nil {
		if vm, ok := v.(map[string]interface{}); ok && vm != nil{
			// try to find param
			delete(vm, param) 
			return nil
		} else {
			// not map, do nothing
			return errors.New("not a map")
		}
	} else {
		// key1 is not found || v is nil
		return errors.New("key1 not found or not a map")
	}
}

//func GetMapValueInt(mp map[string]interface{}, key string) int {
////	fmt.Println(mp, key)
//	if v, pres := mp[key]; pres && v != nil {
////		fmt.Println(v, key, pres)
//		if vint, ok := v.(string); ok {
//			if i, err := strconv.Atoi(vint); err == nil {
//				return i
//			} 
//		} else if vint, ok := v.(float64); ok {
//			return int(vint)
//		} else if vint, ok := v.(int64); ok {
//			return int(vint)
//		}
//		fmt.Println(reflect.TypeOf(v).String())
//	}
//	return 0
//}

func GetMapValueBool(mp map[string]interface{}, key string, dft bool) (bool, error) {
	if v, e := getMapValueAsInt(mp, key); e == nil {
		return v > 0, nil
	} else {
		return dft, e
	}
}

func GetMapValueInt(mp map[string]interface{}, key string, dft int) (int, error) {
	if v, e := getMapValueAsInt(mp, key); e == nil {
		return v, nil 
	} else {
		return dft, e
	}
}


func getMapValueAsInt(mp map[string]interface{}, key string) (int, error) {

	if v, pres := mp[key]; !pres {
		return 0, errors.New("key:" + key + " is not present")
	} else if v == nil {
		return 0, errors.New("key:" + key + " value is nil")
	} else {
		if vint, ok := v.(string); ok {
			if i, err := strconv.Atoi(vint); err == nil {
				return i, nil
			} else {
				return 0, errors.New("key:" + key + " is not a int string")
			}
		} else if vint, ok := v.(float64); ok {
			return int(vint), nil
		} else if vint, ok := v.(int64); ok {
			return int(vint), nil
		} else if vint, ok := v.(int); ok {
			return vint, nil
		} else {
			return 0, errors.New("key:" + key + " unknown type:" + reflect.TypeOf(v).String())
		}
	}
}

func getDifference(a, b []string) (ele []string) {
	ele = make ([]string, 0)
	if a == nil{
		return
	}
	if b == nil {
		ele = append(ele, a...)
		return
	}
	for _, ae := range a {
		found := false
		for _, be := range b {			
			if ae == be {
				found = true				
				break
			}
		} 
		if !found {
			ele = append(ele, ae) 
		} 
	}
	return
} 

func Difference(a, b []string) (newEle, oldEle []string){
	newEle = getDifference(a, b)
	oldEle = getDifference(b, a)
	return
} 

func mapDifference(a, b map[string]interface{}) (ele []string) {
	ele = make ([]string, 0)
	if a == nil{
		return
	}
	if b == nil {
		ele = getMapKeys(a)
		return
	}
	for ae, _ := range a {
		found := false
		for be, _ := range b {			
			if ae == be {
				found = true				
				break
			}
		} 
		if !found {
			ele = append(ele, ae) 
		} 
	}
	return
} 

func getMapKeys(m map[string]interface{}) (ele []string) {
	ele = make ([]string, 0)
	for k, _ := range m {
		ele = append(ele, k)
	}	
	return
}

func MapDifference(a, b map[string]interface{}) (newEle, oldEle []string){
	newEle = mapDifference(a, b)
	oldEle = mapDifference(b, a)
	return
} 

func GetMapKeyValuePair(d map[string]interface{}) ([]string) {
	l := make([]string, 0)
	return AppendMapKeyValuePair(l, d)
}

func AppendMapKeyValuePair(l []string, d map[string]interface{}) ([]string) {

	//fmt.Println(reflect.TypeOf(d).String())

	//	fmt.Println(reflect.TypeOf(d["extra"]).String())
	
//	var err error

	for k, v := range d {
		switch v.(type) {
		case map[string]interface{}:
			l = AppendMapKeyValuePair(l, v.(map[string]interface{}))
		default:
			l = append(l, k+fmt.Sprint(v))
		}
	}
	return l
}


func GetMapValueString2k(mp map[string]interface{}, key1, key2 string) (string, error) {
	
	if v, pres := mp[key1]; pres && v != nil {
		if vm, ok := v.(map[string]interface{}); ok {
			if v1, pres1 := vm[key2]; pres1 && v1 != nil {
				if vm, ok := v1.(string); ok {
					return vm, nil
				} else {
					return "", errors.New(key2 + " not a string")
				}
			} else {
				return "", errors.New(key2 + " not found")
			}
		} else {
			// not map, do nothing
			return "", errors.New(key1 + " not a map")
		}
	} else {
		// key1 is not found
		return "", errors.New(key1 + " not found")
	}
}

func GetMapValueSlice(mp map[string]interface{}, key string) ([]interface{}, error) {
	
	if v, pres := mp[key]; pres && v != nil {
		if vm, ok := v.([]interface{}); ok {
			return vm, nil
		} else {
			// not map, do nothing
			return nil, errors.New("Type Not Found")
		}
	} else {
		// key1 is not found
		return nil, errors.New("Key Not Found")
	}
}

func GetMapValueMap(mp map[string]interface{}, key string) (map[string]interface{}, error) {
	
	if v, pres := mp[key]; pres && v != nil {
		if vm, ok := v.(map[string]interface{}); ok {
			return vm, nil
		} else {
			// not map, do nothing
			return nil, errors.New("Type Not Found")
		}
	} else {
		// key1 is not found
		return nil, errors.New("Key Not Found")
	}
}

func GetIpFromAddr(addr string) string {
	idx := strings.Index(addr, ":")
	if idx > 0 {
		return addr[0:idx]
	}
	return addr
}


func GetMapValueIntSlice(mp map[string]interface{}, key string, sep string) ([]int, error) {
		
	if v, pres := mp[key]; pres && v != nil {
		if vm, ok := v.(string); ok {
			return String2IntSlice(vm, sep)
		} else {
			// not map, do nothing
			return nil, errors.New("Type Not Found")
		}
	} else {
		// key is not found
		return nil, errors.New("Key Not Found")
	}
}
