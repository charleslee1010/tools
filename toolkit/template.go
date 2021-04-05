package toolkit

import (
	"strings"
	"bytes"
	"errors"
	"fmt"
)

/*
 input :
           s   template string, such as:  This $dog and $sheep
           m   parameters map
           replacedWith string to be used to replace the parameters if m is nil
 return:
		 string    after replacement
		 []string  parameters
		 error		              
*/

func TplParse(s string, m map[string]interface{}, replacedWith string) (string, []string, error) {
	var b bytes.Buffer
	plist := make([]string, 0)
	// scan $p
	// replace it with ?
	// put the parameter into p []string
	for true {
		idx := strings.Index(s, "$")
		if idx >= 0 {
			b.WriteString(s[0:idx])
			s = s[idx+1:]
			p, lenp := getNextParam(s)
			if lenp <= 0 {
				// error string
				return "", nil, errors.New("parameter name is not found")
			}
			plist = append(plist, p)

			if m != nil {
				var err error
				replacedWith, err = GetMapValueString(m, p, "")
				if err != nil {
					return "", nil, fmt.Errorf("parameter=%s, err=%v", p, err)
				} 
//				if v, pres := m[p]; pres {
//					if t, ok := v.(string); ok {
//						replacedWith = t
//					}
//				} 
			}
			
			b.WriteString(replacedWith)
			
			s = s[lenp:]
		} else {
			b.WriteString(s[0:])
			break
		}
	}

	return b.String(), plist, nil
}

func getNextParam(s string) (string, int) {

	plen := 0
	for _, v := range s {
		if v >= '0' && v <= '9' ||
			v >= 'a' && v <= 'z' ||
			v >= 'A' && v <= 'Z' ||
			v == '_' {
			plen++
		} else {
			break
		}
	}

	if plen >= 1 {
		return s[0:plen], plen
	} else {
		return "", 0
	}
}


func TplGetParamValues(body map[string]interface{}, plist []string) []interface{} {

	//	log.Debug("input body:%+v, plist=%+v", body, plist)

	v := make([]interface{}, 0)

	for i := 0; i < len(plist); i++ {
		val := body[plist[i]]
		v = append(v, val)
	}

	//	log.Debug("return:%+v", v)

	return v
}

