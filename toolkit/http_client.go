package toolkit

import (
	//	"math/rand"
	"net/http"
	"time"
	//	"srvcfg"
	"crypto/tls"
	"net"
)



func InitHttpClient() *http.Client{
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
		IdleConnTimeout:       10 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
	}

	c := &http.Client{
		Timeout:   10 * time.Second,
		Transport: tr,
	}

	return c
}


func InitHttpClientP(tmo int) *http.Client{
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
		IdleConnTimeout:       10 * time.Second,
		TLSHandshakeTimeout:   5 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
	}

	c := &http.Client{
		Timeout:   time.Duration(tmo) * time.Millisecond,
		Transport: tr,
	}

	return c
}
