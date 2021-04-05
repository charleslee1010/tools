package toolkit

import (
	"strings"
	"os"
	"errors"
	"fmt"
	"path"
	"io/ioutil"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func RenameDuplicateFileName(fn string) error {
	done := false
	
	if pres, _ := PathExists(fn); pres == true {
		// get the filename and suffix
		f, s := SplitSuffix(fn)
		// renmae the file
		for i:= 0; i < 10000; i++ {
			nf := fmt.Sprintf("%s(%d)%s", f, i, s)
			if p, _ := PathExists(nf); p == false {
				//
				os.Rename(fn, nf)
				done = true 
				break
			}			
		}		
	} else {
		done = true
	}
	
	if done {
		return nil
	} else {
		return errors.New("can not rename")
	}
}

func CreateFile(fn string) (*os.File, error) {
	err := RenameDuplicateFileName(fn)
	
	if err == nil {
		return os.Create(fn)
	} else {
		return nil, err
	}
}

func SplitSuffix(fn string) (string, string) {
	idx := strings.LastIndexByte(fn, '.')
	if idx < 0 {
		return fn, ""
	}
	return fn[0:idx], fn[idx:]
}
	
func BatchRename(from, to string, rn bool) error {
	bf := path.Base(from)
	bt := path.Base(to)
	
	df := path.Dir(from)
	dt := path.Dir(to)
	
	
	fmt.Println(df, ",", bf)
	fmt.Println(dt, ",", bt)
	
	if f, err := ioutil.ReadDir(df); err != nil {
		return err
	} else {
		for _, v := range f {
			if v.IsDir() {
				continue
			}
			fmt.Println(v.Name())
			if r, e := MatchAndReplace(v.Name(), bf, bt); e == nil && r != bt {
				fname := path.Join(df, v.Name())
				tname := path.Join(dt,r)
				
				fmt.Printf("rename from %s to %s\n", fname, tname)
				
				if rn {
					if e := os.Rename(fname, tname); e != nil {
						fmt.Printf("rename failed, err=\n", e)
					}
				}
				
			} else {
				fmt.Printf("rename error %v or invalid target:%s,%s\n", e, r, bt)
			}
		}
		return nil
	}
}	
