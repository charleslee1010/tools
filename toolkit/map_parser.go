package toolkit

import (
	"strconv"
	"errors"
	"reflect"
)

var ERR_KEY_NOT_FOUND error = errors.New("key not found")


type MapParser struct {
	root       interface{}
	m          interface{} 
	lastErr    error                 // last error reported
}

func NewMapParser(mp map[string]interface{}) *MapParser {
	return &MapParser{ m: mp, root: mp}
}

func IsKeyNotFoundErr(err error) bool {
	
	return err == ERR_KEY_NOT_FOUND
}


func (this*MapParser) Restart() *MapParser  {
	this.m = this.root
	this.lastErr = nil
	return this
}
func (this*MapParser) Reset() *MapParser  {
	this.m = this.root
	this.lastErr = nil
	return this
}
func (this*MapParser) String() (string, error) {
	if v, ok := this.m.(string); ok {
		return v, nil
	} else if v, ok := this.m.(float64); ok {
		return strconv.FormatFloat(v, 'f', -1, 64), nil
	} else if v, ok := this.m.(int64); ok {
		return strconv.FormatInt(v, 10), nil
	} else {
		return "", errors.New(reflect.TypeOf(this.m).String())		
	}
}

func (this*MapParser) Int() (int, error) {
	if v, ok := this.m.(float64); ok {
		return int(v), nil
	} else if v, ok := this.m.(int64); ok {
		return int(v), nil
	} else if v, ok := this.m.(int); ok {
		return v, nil
	} else if v, ok := this.m.(string); ok {
		return strconv.Atoi(v)
	} else {
		return 0, errors.New(reflect.TypeOf(this.m).String())		
	}
}

func (this*MapParser) Int64() (int64, error) {
	if v, ok := this.m.(float64); ok {
		return int64(v), nil
	} else if v, ok := this.m.(int64); ok {
		return int64(v), nil
	} else if v, ok := this.m.(string); ok {
		v, err:= strconv.Atoi(v)
		return int64(v), err
	} else {
		return 0, errors.New(reflect.TypeOf(this.m).String())		
	}
}

func (this*MapParser) Map(key string) *MapParser {
	if this.lastErr != nil {
		return this
	}
	
	if v, ok := this.m.(map[string]interface{}); ok {
		if v1, pres := v[key]; pres {
			this.m = v1
		} else {
			this.lastErr = ERR_KEY_NOT_FOUND   //errors.New("key not found")
		}
	} else {
		this.lastErr = errors.New(reflect.TypeOf(this.m).String())
	}
	return this
}



func (this*MapParser) Slice(id int) *MapParser {
	if this.lastErr != nil {
		return this
	}
	
	if v, ok := this.m.([]interface{}); ok {
		if id < len(v) {
			this.m = v[id]
		} else {
			this.lastErr = errors.New("id out of range")
		}
	} else {
		this.lastErr = errors.New(reflect.TypeOf(this.m).String())
	}
	return this
}


func (this*MapParser) StringSliceAt(id int) *MapParser {
	if this.lastErr != nil {
		return this
	}
	
	if v, ok := this.m.([]string); ok {
		if id < len(v) {
			this.m = v[id]
		} else {
			this.lastErr = errors.New("id out of range")
		}
	} else {
		this.lastErr = errors.New(reflect.TypeOf(this.m).String())
	}
	return this
}

func (this*MapParser) StringSlice() ([]string, error) {
	if this.lastErr != nil {
		return nil, this.lastErr
	}
	
	if v, ok := this.m.([]string); ok {
		return v, nil
//		if id < len(v) {
//			this.m = v[id]
//		} else {
//			this.lastErr = errors.New("id out of range")
//		}
	} else {
		this.lastErr = errors.New(reflect.TypeOf(this.m).String())
	}
	return nil, nil
}

func (this*MapParser) NewMap(key string) (*MapParser) {
	if this.lastErr != nil {
		return this
	}
	
	if v, ok := this.m.(map[string]interface{}); ok {
		if v1, pres := v[key]; pres {
			if _, ok := v1.(map[string]interface{}); ok {
				// already exist, do nothing
				this.m = v1
			} else {
				// node exist other than map
				this.lastErr = errors.New("key value is not a map")
			}
		} else {
			m := make(map[string]interface{})
			v[key] = m
			this.m = m
		}
	} else {
		this.lastErr = errors.New("cannot create map for non-map structure")
	}
	return this
}


func (this*MapParser) SetMapKeyValue(k string, newv interface{}) (error) {
	if this.lastErr != nil {
		return this.lastErr
	}
	
	if v, ok := this.m.(map[string]interface{}); ok {
		v[k] = newv
		return nil
	} else {
		// not a map, cannot set value
		return errors.New("key value is not a map")
	}
}

