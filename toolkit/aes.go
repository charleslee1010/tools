package toolkit

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	
)

// 填充数据
func padding(src []byte, blockSize int) []byte {
	padNum := blockSize - len(src)%blockSize
	pad := bytes.Repeat([]byte{byte(padNum)}, padNum)
	return append(src, pad...)
}

// 去掉填充数据
func unpadding(src []byte) []byte {
	n := len(src)
	unPadNum := int(src[n-1])
	return src[:n-unPadNum]
}

// 加密
func encryptAES(src []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	src = padding(src, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, key)
	blockMode.CryptBlocks(src, src)
	return src, nil
}

// 解密
func decryptAES(src []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode := cipher.NewCBCDecrypter(block, key)
	blockMode.CryptBlocks(src, src)
	src = unpadding(src)
	return src, nil
}


var	key = []byte("hgAedc4a87604321")


func AesEncrypt(s string) string {
	x1, err := encryptAES([]byte(s), key)
	if err != nil {
		return ""
	}

	return hex.EncodeToString(x1)
}

func AesDecrypt(s string) string {
	if b, err := hex.DecodeString(s); err == nil {	
		if x1, err := decryptAES(b, key); err == nil {
			return string(x1)
		}
	}

	return ""
}




