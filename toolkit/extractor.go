package toolkit

import (
//	"strings"
	"strconv"
	"errors"
//	"fmt"
	"reflect"

)


type Extractor struct {
	m          interface{} 
	Error    error                 // last error reported
}

func NewExtractor(mp map[string]interface{}) *Extractor {
	return &Extractor{ m: mp, }
}

func (this*Extractor) String() (string, error) {
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

func (this*Extractor) Int() (int, error) {
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

func (this*Extractor) Int64() (int64, error) {
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

func (this*Extractor) Map(key string) *Extractor {
	if this.Error != nil {
		return this
	}
	
	if v, ok := this.m.(map[string]interface{}); ok {
		if v1, pres := v[key]; pres {
			this.m = v1
		} else {
			this.Error = errors.New("key not found")
		}
	} else {
		this.Error = errors.New(reflect.TypeOf(this.m).String())
	}
	return this
}



func (this*Extractor) Slice(id int) *Extractor {
	if this.Error != nil {
		return this
	}
	
	if v, ok := this.m.([]interface{}); ok {
		if id < len(v) {
			this.m = v[id]
		} else {
			this.Error = errors.New("id out of range")
		}
	} else {
		this.Error = errors.New(reflect.TypeOf(this.m).String())
	}
	return this
}

func (this*Extractor) GetSlice() [] interface{} {
	if this.Error != nil {
		return nil
	}
	
	if v, ok := this.m.([]interface{}); ok {
		return v
	} else {
		return nil
	}
}


func (this*Extractor) StringSlice(id int) *Extractor {
	if this.Error != nil {
		return this
	}
	
	if v, ok := this.m.([]string); ok {
		if id < len(v) {
			this.m = v[id]
		} else {
			this.Error = errors.New("id out of range")
		}
	} else {
		this.Error = errors.New(reflect.TypeOf(this.m).String())
	}
	return this
}

