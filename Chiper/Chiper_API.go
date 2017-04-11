package Chiper

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

var decrypted string
var privateKey, publicKey []byte

// 加密
func RsaEncrypt(origData []byte) ([]byte, error) {
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
func RsaDecrypt(ciphertext []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

//RSA公钥私钥产生
func GenRsaKey(bits int) error {
	// 生成私钥文件
	//var privateKey *rsa.PrivateKey
	//privateKey =new(rsa.PrivateKey)
	PKEY, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	derStream := x509.MarshalPKCS1PrivateKey(PKEY)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	file, err := os.Create("d:\\private.pem")
	if err != nil {
		return err
	}
	privateKey = pem.EncodeToMemory(block)
	fmt.Println(privateKey)

	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	// 生成公钥文件
	PUKEY := &PKEY.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(PUKEY)
	if err != nil {
		return err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	file, err = os.Create("d:\\public.pem")
	if err != nil {
		return err
	}
	publicKey = pem.EncodeToMemory(block)
	err = pem.Encode(file, block)
	if err != nil {
		return err
	}
	return nil
}

func load_keys() {
	var err error
	// flag.StringVar(&decrypted, "d", "", "加密过的数据")
	// flag.Parse()
	publicKey, err = ioutil.ReadFile("d:\\public.pem")
	if err != nil {
		os.Exit(-1)
	}
	privateKey, err = ioutil.ReadFile("d:\\private.pem")
	if err != nil {
		os.Exit(-1)
	}
}
