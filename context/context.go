package context

import (
	"encoding/json"
	"github.com/charles/mylog"
	"io/ioutil"
	"net/http"
	"time"
	uuid "github.com/satori/go.uuid"
)

type Context struct {
	W       http.ResponseWriter
	R       *http.Request
	BodyStr string
	Sid     string
	Logger  *mylog.MyLogger
	StartTime	time.Time
}
//

func NewContext(w http.ResponseWriter, r *http.Request, logger *mylog.MyLogger) *Context {

	return &Context{
		W:      w,
		R:      r,
		Sid:    uuid.Must(uuid.NewV4()).String(),
		Logger: logger,
		StartTime:     time.Now(),
	}
}

func (this *Context) Init(w http.ResponseWriter, r *http.Request, logger *mylog.MyLogger) {
	this.W = w
	this.R = r
	this.Sid = uuid.Must(uuid.NewV4()).String()
	this.Logger = logger
	this.StartTime = time.Now()
}

func (this *Context) SendJson(code int, v interface{}) {

	if b, err := json.Marshal(v); err == nil {
		this.Logger.Trace("RSP -> APP sid=%s, code=%d, resp=%s", this.Sid, code, string(b))
		this.W.WriteHeader(code)
		this.W.Write(b)
	} else {
		m := "malformed body"
		this.Logger.Trace("RSP -> APP sid=%s, code=%d, resp=%s", this.Sid, code, m)
		this.W.WriteHeader(code)
		this.W.Write([]byte(m))
	}
}

func (this *Context) SendBinary(code int, b []byte) {

	this.Logger.Trace("RSP -> APP sid=%s, code=%d, resp=%s", this.Sid, code, string(b))
	this.W.WriteHeader(code)
	this.W.Write(b)
}

func (this *Context) SendStatusCode(code int) {

	this.Logger.Trace("RSP -> APP sid=%s, code=%d", this.Sid, code)
	this.W.WriteHeader(code)
}

func (this *Context) ReadJson(v interface{}) error {

	s, _ := ioutil.ReadAll(this.R.Body) //把  body 内容读入字符串 s
	this.BodyStr = string(s)

	str := this.BodyStr
	if len(str) > 250 {
		str = str[0:250]
	}

	// print Logger
	this.Logger.Trace("REQ <- APP, sid=%s %s %s %s %s body=%s",
		this.Sid, this.R.Method, this.R.RemoteAddr, this.R.Host, this.R.URL.String(), str)

	if len(s) != 0 {
		if err := json.Unmarshal(s, v); err != nil {
			this.Logger.Info("Sid=%s, malformed message body, err=%+v", this.Sid, err)
			return err
		}
	}
	return nil
}

func (this *Context) ReadBinary() ([]byte, error) {

	b, err := ioutil.ReadAll(this.R.Body) //把  body 内容读入字符串 b

	str := string(b)
	if len(str) > 250 {
		str = str[0:250]
	}
	// print Logger
	this.Logger.Trace("REQ <- APP, sid=%s %s %s %s %s body=%s",
		this.Sid, this.R.Method, this.R.RemoteAddr, this.R.Host, this.R.URL.String(), str)

	return b, err
}
