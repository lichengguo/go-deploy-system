package aespwd

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"go-deploy-system-server/utils"
)

// 高级加密标准（Adevanced Encryption Standard ,AES）

// 16,24,32位字符串的话，分别对应AES-128，AES-192，AES-256 加密方法
// key不能泄露
var pwdKey = []byte(utils.ServerGitKey)

// pkcs7Padding PKCS7 填充模式
func pkcs7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize

	// Repeat()函数的功能是把切片[]byte{byte(padding)}复制padding个，然后合并成新的字节切片返回
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(ciphertext, padtext...)
}

// pkcs7UnPadding 填充的反向操作，删除填充字符串
func pkcs7UnPadding(origData []byte) ([]byte, error) {
	// 获取数据长度
	length := len(origData)
	if length == 0 {
		return nil, errors.New("加密字符串错误")
	}

	// 获取填充字符串长度
	unpadding := int(origData[length-1])

	// 截取切片，删除填充字节，并且返回明文
	return origData[:(length - unpadding)], nil
}

// aesEcrypt 实现加密
func aesEcrypt(origData []byte, key []byte) ([]byte, error) {
	// 创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 获取块的大小
	blockSize := block.BlockSize()

	// 对数据进行填充，让数据长度满足需求
	origData = pkcs7Padding(origData, blockSize)

	// 采用AES加密方法中CBC加密模式
	blocMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	crypted := make([]byte, len(origData))

	// 执行加密
	blocMode.CryptBlocks(crypted, origData)

	return crypted, nil
}

// aesDeCrypt 实现解密
func aesDeCrypt(cypted []byte, key []byte) ([]byte, error) {
	// 创建加密算法实例
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 获取块大小
	blockSize := block.BlockSize()

	// 创建加密客户端实例
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(cypted))

	// 这个函数也可以用来解密
	blockMode.CryptBlocks(origData, cypted)

	// 去除填充字符串
	origData, err = pkcs7UnPadding(origData)
	if err != nil {
		return nil, err
	}

	return origData, err
}

// EnPwdCode 加密base64
func EnPwdCode(pwdStr string) (string, error) {
	var pwd []byte
	pwd = []byte(pwdStr)

	result, err := aesEcrypt(pwd, pwdKey)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(result), err
}

// DePwdCode 解密base64
func DePwdCode(pwd string) (string, error) {
	// 解密base64字符串
	pwdByte, err := base64.StdEncoding.DecodeString(pwd)
	if err != nil {
		return "", err
	}

	// 执行AES解密
	pwsBtye, _ := aesDeCrypt(pwdByte, pwdKey)

	return string(pwsBtye), nil
}
