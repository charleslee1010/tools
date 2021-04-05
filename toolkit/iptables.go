package toolkit

import (
//	"strings"
//	"strconv"
	"errors"
	"fmt"
//	"reflect"
	"os"
//	"strings"
//	"path/filepath"
	"os/exec"
	
	"bytes"

)
type IpTablesStruct struct {
	fn        string           // file name
	
	header, footer string		// header filename, footer file name
	
	template	string		// iptables template line
	
	// do restore only if it is "true"
	restore     string
	
	Total		int
}


func NewIpTables(fn string, header, footer, template, restore string) (*IpTablesStruct, error){
	ips := &IpTablesStruct{
		fn:fn,
		header:header,
		footer:footer,
		template:template,
		restore:restore,
	}
	
	// check if fn exists, open for read and write
	if f, err := os.OpenFile(header, os.O_RDONLY, 0755); err != nil {
		return nil, err
	} else {
		f.Close()		
	}

	if f, err := os.OpenFile(footer, os.O_RDONLY, 0755); err != nil {
		return nil, err
	} else {
		f.Close()		
	}

	if f, err := os.OpenFile(fn, os.O_CREATE|os.O_TRUNC, 0755); err != nil {
		return nil, err
	} else {
		f.Close()		
	}
    
    return ips, nil
}

const ShellToUse = "/bin/bash"

func Shellout(command string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}

func RunShell(command string) (string, string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command(ShellToUse, "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return stdout.String(), stderr.String(), err
}
//
//// addr format :  ip:port
//func (this*IpTablesStruct)AddInput(addr string) error {
//	t := strings.Split(addr, ":")
//	if t == nil || len(t) > 2 {
//		return errors.New("invalid address")
//	}
//	// t[0] ip addr   t[1] port
////	port := ""
////	if len(t) == 2 && t[1] != "*" {
////		port = t[1]
////	}
//	
//	//sed
//	// sed -i '/特定字符串/a\server 127.127.1.0' 
//	cmd := "-A INPUT -s " + t[0] + "/32 -p udp --dport 15063 -j ACCEPT"
//
//	sed := `sed -i '/clientaddress/a\` + cmd + `' ` + this.fn
//	if err, _, _:= Shellout(sed); err != nil {
//		return errors.New(sed + " err:" + err.Error())
//	}
////	fmt.Println("ok: " + sed)
//	// iptables-restore
//	rst := `sudo iptables-restore <` + this.fn
//	if err, _, _ := Shellout(rst); err!=nil {
//		return errors.New(rst + " err:" + err.Error())
//	}
////	fmt.Println("ok: " + rst)
//	return nil
//}
//
//// addr format :  ip:port
//func (this*IpTablesStruct)DelInput(addr string) error {
//	t := strings.Split(addr, ":")
//	if t == nil || len(t) > 2 {
//		return errors.New("invalid address")
//	}
//	// t[0] ip addr   t[1] port
////	port := ""
////	if len(t) == 2 && t[1] != "*" {
////		port = t[1]
////	}
////	
//	//sed
//	// sed -i '/特定字符串/a\server 127.127.1.0' 
//	// sed -i "/pattern/d"  log.txt
//	ipt := strings.Split(t[0], ".")
//	if len(ipt) != 4 {
//		return errors.New("invalid ip")
//	}
//	sed := `sed -i '/` + ipt[0]+`\.`+ ipt[1]+`\.`+ ipt[2]+`\.`+ ipt[3] +`/d' ` + this.fn
//	if err, _, _:= Shellout(sed); err != nil {
//		return errors.New(sed + " err:" + err.Error())
//	}
//	
////	fmt.Println("ok: " + sed)
//	
//	// iptables-restore
//	rst := `sudo iptables-restore <` + this.fn
//	if err, _, _ := Shellout(rst); err!=nil {
//		return errors.New(rst + " err:" + err.Error())
//	}
////	fmt.Println("ok: " + rst)
//	
//	return nil
//}
//

// addr format :  ip:port
func (this*IpTablesStruct)Restore(ips []map[string]interface{}) error {

	if err := this.Generate(ips); err != nil {
		return err
	} else if this.restore == "true" {	
		// iptables-restore
		rst := `sudo /sbin/iptables-restore <` + this.fn
		if err, s1, s2 := Shellout(rst); err!=nil {			
			return fmt.Errorf("cmd=%s, err=%v, s1=%v, s2=%v", rst, err, s1, s2)
		}
	}
	
	return nil
}

// addr format :  ip:port
func (this*IpTablesStruct)Generate(ips []map[string]interface{}) error {
	this.Total = 0
	// check if fn exists, open for read and write
	fout, err := os.OpenFile(this.fn, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer fout.Close()
	
	// read from txt1 and write into fout
	
	if err := copyFile(fout, this.header); err != nil {
        return err
    }

	// write m into fout
	
	if err := this.writeIpRules(fout, ips); err != nil {
        return err
    }
	
	// read from txt1 and write into fout
	if err := copyFile(fout, this.footer); err != nil {
        return err
    }
	
	// check number of lines of fout
	if fi, err := fout.Stat(); err == nil && fi.Size() < 100 {
		return fmt.Errorf("invalid iptables file length:%d", fi.Size())
	}
    
	return nil
}


func (this*IpTablesStruct)writeIpRules(fout*os.File, r []map[string]interface{}) error {

    fout.WriteString("\n")
	
	for _, v := range r {
		// ip must be present
		if ip, err := GetMapValueString(v, "ip", ""); err != nil || ip == "" {
			return errors.New("ip is empty")
		}
		
		if a, _, err := TplParse(this.template, v, ""); err != nil {
			return err
		} else {
			fout.WriteString(a)
			fout.WriteString("\n")
		}		
	}
	this.Total = len(r)
	return nil
}
const LEN_BUF = 1024
func copyFile(fout*os.File, fn string) error {
	f, err := os.OpenFile(fn, os.O_RDONLY, 0755)
    if err != nil {
        return err
    }
    
    defer f.Close()
    
    b := make([]byte, LEN_BUF)
	for true {
		if n, _ := f.Read(b); n > 0 {
			fout.Write(b[0:n])
//			fmt.Println("write ", n, " bytes")
		} else {
			break
		}
	}
	return nil
}
