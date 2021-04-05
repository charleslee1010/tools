package toolkit

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"testing"
)
	
	
func TestGetProtocolIpPort(t *testing.T) {
	d0 := "http://47.110.238.181:12401"
	d1 := "/V2/20200907155945/ZHVIP_17612598364_13979972784_17158619711_20200907155945_425f55e871359451.mp3"
	r0, r1 := GetProtocolIpPort(d0+d1)
	if r0 != d0 || r1 != d1{
		t.Errorf("error1:%s, %s", r0, r1)
	}


	d0 = "https://47.110.238.181:12401"
	d1 = "/V2/20200907155945/ZHVIP_17612598364_13979972784_17158619711_20200907155945_425f55e871359451.mp3"
	r0, r1 = GetProtocolIpPort(d0+d1)
	if r0 != d0 || r1 != d1{
		t.Errorf("error2:%s, %s", r0, r1)
	}

	d0 = "https://47.110.238.181:12401"
	r0,_ = GetProtocolIpPort(d0)
	if r0 != "" {
		t.Errorf("error3:%s", r0)
	}

	d0 = "http:/47.110.238.181:12401/"
	r0,_ = GetProtocolIpPort(d0)
	if r0 != "" {
		t.Errorf("error4:%s", r0)
	}		
}



func TestAes(t *testing.T) {
	//d := []byte("http://47.110.238.181:12401/V2/20200907155945/ZHVIP_17612598364_13979972784_17158619711_20200907155945_425f55e871359451.mp3")
	d := []byte("47.110.238.181:12401")
	key := []byte("hgfedcba87654321")
	fmt.Println("加密前:", string(d))
	x1, err := encryptAES(d, key)
	if err != nil {
		t.Error("enc error")
	}

	fmt.Println("密文(hex)：", hex.EncodeToString(x1))
	fmt.Println("密文(base64)：", base64.StdEncoding.EncodeToString(x1))

	fmt.Println("加密后:", string(x1))
	x2, err := decryptAES(x1, key)
	if err != nil {
		t.Error("dec error")
	}
	fmt.Println("解密后:", string(x2))

	if string(x2) != string(d) {
		t.Error("mismatch")
	}
}



func TestZlib(t *testing.T) {
//    buff := []byte{120, 156, 202, 72, 205, 201, 201, 215, 81, 40, 207,
//        47, 202, 73, 225, 2, 4, 0, 0, 255, 255, 33, 231, 4, 147}
//    b := bytes.NewReader(buff)
//    r, err := zlib.NewReader(b)
//    if err != nil {
//        panic(err)
//    }
//    io.Copy(os.Stdout, r)
//    r.Close()
	s := "47.110.238.181:12401"
    zip := DoZlibCompress([]byte(s))
    fmt.Println(len(zip), zip)
    s0 := string(DoZlibUnCompress(zip))
	if s0 != s {
		t.Error("failed")
	}
}

func TestAes1(t *testing.T) {
	//d := []byte("http://47.110.238.181:12401/V2/20200907155945/ZHVIP_17612598364_13979972784_17158619711_20200907155945_425f55e871359451.mp3")
	s := "47.110.238.181:12401"
	
	d := AesEncrypt(s)
	if d == "" {
		t.Error("enc error")
		return
	}
	
	s0 := AesDecrypt(d)
	if s0 == "" {
		t.Error("dec error")
		return
	}
	
	if s != s0 {
		t.Error("mismatch")
	}
}

func TestMapParser(t *testing.T) {
	m := make(map[string]interface{})
	
	mp := NewMapParser(m).NewMap("key1").NewMap("key2").NewMap("key2").NewMap("key3")
	mp.SetMapKeyValue("key4", 101)
	fmt.Printf("%+v\n", mp)
	
	// overwrite
	NewMapParser(m).NewMap("key1").NewMap("key2").NewMap("key2").SetMapKeyValue("key3", "string")
	fmt.Printf("%+v\n", m)
	
//	fmt.Printf("%+v\n", mp.NewMap("key1"))
	
//	.SetMapKeyValue("key4", 101)
}
	

