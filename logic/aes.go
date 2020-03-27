package logic

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
)

var publicKey = []byte(`
-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQDJXWWPKYpJTLTeBz78u6NIx3vf
aGjAXdO00JAJHhKjKbVeUtXZ9xIvOY86EYTtJ4rW0cd+ERHVjSOXLPo2OBfIqY33
XfVSBbbJleXAvCRAXw5idxEQ9ZKESFGogoKHlKzKia1OAPczJxvMinhmHMhNF1fY
T+n5sfbVoLohgdex2wIDAQAB
-----END PUBLIC KEY-----
`)

var privateKey = []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQDJXWWPKYpJTLTeBz78u6NIx3vfaGjAXdO00JAJHhKjKbVeUtXZ
9xIvOY86EYTtJ4rW0cd+ERHVjSOXLPo2OBfIqY33XfVSBbbJleXAvCRAXw5idxEQ
9ZKESFGogoKHlKzKia1OAPczJxvMinhmHMhNF1fYT+n5sfbVoLohgdex2wIDAQAB
AoGAT+E6/pXOA9HoFgPt2rhcx+xKmY+DrnwKFbp+yP8jCZLsHeTibLr0fcNpq/Fz
N9jt3NYPO1VuK7b3nWr8PzH1TLK9Q3ykC6lrxA65QjL9rylSqsqbn8vpU8BlQlq1
KjounkhzsZ4oSPNA7WUlIChpTk3Dl2+n36avoZndXKlu0pECQQDmyDbCHXHHmh5R
T7l9bpKmbnB8ge0Ug3J0iZ4DiF8nBIX2SxUYcHTcQXiFWfgt6YdIRaagQhZ59Z1Y
xHEfytWFAkEA315HX+FTdOzvsvoGQf/sPTmYZhBJy1hflMplmZaYsId5G7ln+l6x
48vYy4dsjx/GY59jEOQZMaugrtm39RHX3wJBAL/V29a9/Q9raBo1CD5gxJx+HxkQ
M0+i+Ggw4N2U5WuckfKqdO2sxSc1cQaARBF+FosYAqsiZGaaqWHZYSOJSrUCQQCe
D/EiACk2jJPyassS2S8rBB672rrdkmPQvoi27sKN6M/itojFu0zWjeGT5PkFLs8M
oDVSNpc9dt313Us3uLCxAkAwtDUUB5WvPYz5X6HUSX/yDw8jPDd4KsLwDF+Nmfbu
Juh9d1hIpVx96Mn2S12wuvOJBCiKH3tXcvjBUGh0Sv80
-----END RSA PRIVATE KEY-----
`)

// 加密
func rsaEncrypt(origData []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// 解密
func rsaDecrypt(cipherText []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, cipherText)
}

//加密再次封装
func encrypt(origData string) (string, error) {
	data, err := rsaEncrypt([]byte(origData))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

//解密再次封装
func decrypt(decrypted string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(decrypted)
	if err != nil {
		return "", err
	}
	origData, err := rsaDecrypt(data)
	if err != nil {
		return "", err
	}
	return string(origData), nil
}

//分段加密
func splitEncrypt(origData string) (string, error) {
	sLen := len(origData)
	var jmSlice []string
	for i := 0; i <= (sLen / 80); i++ {
		var ns string
		if 80*(i+1) > sLen {
			ns = origData[80*i:]
		} else {
			ns = origData[80*i : 80*(i+1)]
		}
		js1, err := encrypt(ns)
		if err != nil {
			return "", err
		}
		jmSlice = append(jmSlice, js1)
	}
	//序列化
	data, err := json.Marshal(jmSlice)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
