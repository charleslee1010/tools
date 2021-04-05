package toolkit

import (
	"strings"
	"sync"
	"strconv"
	"errors"
)

type Properties struct {
	L		sync.RWMutex		
	Pmap map[string] []string // properties map
	
	FileName string
}

// clear the value
func (prop*Properties) Clear() {
	prop.Pmap = nil
}


const (
	STICKY_TAG = `\`
)

func (this*Properties)ReloadFromFile(fn string) (error) {
	this.L.Lock()
	defer this.L.Unlock()
	
	if fn != "" {
		this.FileName = fn
	}
	
	this.Pmap = make(map[string][]string)
	
	// the last line NOT yet processed
	lastLine := ""
	// the last line contains the sticky tag \
	lastSticky := false

	err := ReadTextFile(this.FileName, func (lno int, line string) bool{
		// now we have line
		
//		fmt.Println("got line:", line)

		// check if it is comment line
		line = removeComments(line)
		//
		if len(line) > 0 && line[len(line)-1:] == STICKY_TAG {
			// this line contains the sticky tag, append it to the last line
			lastLine = lastLine + line[0:len(line)-1]
			lastSticky = true
			// do nothing and return
			return true 
		} else {
			// current line is not sticky
			if lastSticky {
				// the last line is sticky
				line = lastLine + line
				lastLine = ""
				lastSticky = false
			} else {
				// the last is not sticky
			}
		}
		
		// now we have a line, process it		
		
		kIdx := strings.Index(line, "=")

		if kIdx <= 0 {
//			fmt.Println("ignored")
			return true
			// invalid line， ignore it
		}

		this.Add(line[0:kIdx], line[kIdx+1:])
		
		return true
	})
	
	
	// if the last line is sticky, treat it as error
	if lastSticky {
		err = errors.New("sticky tag not end")
	}
	return err
}


func (this*Properties)Add(k, v string) {
		
		key := strings.TrimSpace(k)
		value := strings.TrimSpace(v)
		
		if s, ok := this.Pmap[key]; ok {
			this.Pmap[key] = append(s, value)
		} else {
			 this.Pmap[key] = []string {value}
		}
}


// two consecutive ##, is treated as one valid #
func removeComments(line string) string {	
	ret := ""
	left := line
	
	for true { 	
		pidx := strings.Index(left, "#")
		if pidx >= 0 {
			// save the valid
			if pidx+1  == len(left) {
				// it is the last char
				if pidx > 0 {
					ret = ret + left[0:pidx]
				}
				break
			} else {
				// more char			
				ret = ret + left[0:pidx]
				// check the next char
				if left[pidx+1:pidx+2] == "#" {
					// consecutive #
					ret = ret + "#"
					if pidx+2 == len(left) {
						break
					} else {
						left = left[pidx+2:]
					}
				} else {
					// no consecutive #
					break
				}
			}
			
		} else {
			// no more # found
			ret = ret + left
			break
		}
	}
	return strings.TrimSpace(ret)
}

//func NewProperties(fname string) (*Properties, error) {
//	prop := &Properties{
//		Pmap :make(map[string][]string),
//	}
//
//	err := ReadTextFile(fname, func (lno int, line string) bool{
//		// now we have line
//		
////		fmt.Println("got line:", line)
//
//		// check if it is comment line
//		line = strings.TrimSpace(line)
//		pidx := strings.Index(line, "#")
//		if pidx >= 0 {
//			line = line[0:pidx]
//		}
//		//
//		
//		kIdx := strings.Index(line, "=")
//
//		if kIdx <= 0 {
////			fmt.Println("ignored")
//			return true
//			// invalid line， ignore it
//		}
//
//		key := strings.TrimSpace(line[0:kIdx])
//		value := strings.TrimSpace(line[kIdx+1:])
//		
//		if s, ok := prop.Pmap[key]; ok {
//			// append the string to
//			prop.Pmap[key] = append(s, value)
//		} else {
//			 prop.Pmap[key] = []string {value}
//		}
//		
//		return true
//	})
//	
//	return prop, err
//}

func (prop*Properties) Get(key string) (string) {
	prop.L.RLock()
	defer prop.L.RUnlock()
	
	if v, ok := prop.Pmap[key]; ok && len(v) >= 1 {
		return v[0]
	}
	return ""
}

func (prop*Properties) GetPropertyIntRange(key string, dft int, min, max int) (int) {
	prop.L.RLock()
	defer prop.L.RUnlock()

	v := prop.Get(key)
	
	if v == "" {
		return dft
	}
	
	if ret, err:=strconv.Atoi(v); err != nil || ret > max || ret < min {
		return dft
	} else {
		return ret 
	}
}

func (prop*Properties) GetPropertyInt(key string, dft int) (int) {
	prop.L.RLock()
	defer prop.L.RUnlock()

	v := prop.Get(key)
	
	if v == "" {
		return dft
	}
	
	if ret, err:=strconv.Atoi(v); err != nil {
		return dft
	} else {
		return ret 
	}
}

func (prop*Properties) GetPropertyString(key string, dft string) (string) {
	prop.L.RLock()
	defer prop.L.RUnlock()

	v := prop.Get(key)

	if v == "" {
		return dft
	}
	return v
}


//
// key=abc,dfdf,sdfa,dsf
// return ["abc","dfdf","sdfa","dsf"]
//

func (prop*Properties) GetPropertySlice(key string, sep string) ([]string) {
	prop.L.RLock()
	defer prop.L.RUnlock()

	v := prop.Get(key)

	if v == "" {
		return []string{}
	}
	
	t := strings.Split(v, sep)
	
	for i := 0; i < len(t); i++ {
		t[i] = strings.TrimSpace(t[i])
	}
	
	return t
}

func (prop*Properties) GetPropertyStringList(key string) ([]string) {

	prop.L.RLock()
	defer prop.L.RUnlock()
	
	if v, ok := prop.Pmap[key]; ok {
		return v
	}
	return make([]string, 0)
}

//
// key=k1, v1
// key=k2, v2
// key=k3, v3
// return map[] {"k1":"v1","k2":"v2","k3":"v3"}
//
func (prop*Properties) GetPropertyKeyMap(key string) (m map[string]string) {
	prop.L.RLock()
	defer prop.L.RUnlock()

	m = make(map[string]string)

	if v, ok := prop.Pmap[key]; ok {
		// now we have v []string			
		for _, s := range v {
			// find the first ,
			idx := strings.Index(s, ",")
			if idx > 0 {
				m[strings.TrimSpace(s[0:idx])] = strings.TrimSpace(s[idx+1:])
			}			
		}
	}
	return
}

//
// key=1,2,545,656
// return:     [1,2,545,656]
//
func (prop*Properties) GetPropertySliceInt(key string, sep string) ([]int) {
	prop.L.RLock()
	defer prop.L.RUnlock()

	v := prop.Get(key)

	if v == "" {
		return nil
	}
	
	t := strings.Split(v, sep)
	
	intSlice := make([]int, 0)
	
	for i := 0; i < len(t); i++ {
		if iv, err := strconv.Atoi(t[i]); err == nil {
			intSlice = append(intSlice, iv)  
		} else {
			return nil
		}
	}
	
	return intSlice
}
