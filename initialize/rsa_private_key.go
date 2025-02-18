package initialize

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

var RSAPrivateKey *rsa.PrivateKey

func RsaDecryptInit(path string) (err error) {
	filePath := path + "/rsa_key/private_key.pem"
	key, err := os.ReadFile(filePath)
	if err != nil {
		return errors.New("加载私钥错误1：" + err.Error())
	}
	block, _ := pem.Decode(key)
	if block == nil {
		return errors.New("加载私钥错误2：")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return errors.New("加载私钥错误3：" + err.Error())
	}
	RSAPrivateKey = privateKey
	return err
}

func DecryptPassword(encryptedPassword string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encryptedPassword)
	if err != nil {
		return nil, fmt.Errorf("解码密文失败: %v", err)
	}

	decrypted, err := rsa.DecryptPKCS1v15(rand.Reader, RSAPrivateKey, ciphertext)
	if err != nil {
		return nil, fmt.Errorf("解密失败: %v", err)
	}

	return decrypted, nil
}
