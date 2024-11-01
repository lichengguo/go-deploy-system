package md5

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"math/rand"
	"time"
)

// 生成随机数字符串
func randomStr() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := r.Intn(1024)

	return string(num)
}

// 文件名md5后存储，避免同名覆盖旧文件
func StringToMD5(fileName string) (newFileName string) {
	h := md5.New()
	_, _ = io.WriteString(h, randomStr())
	_, _ = io.WriteString(h, fileName)
	newFileName = hex.EncodeToString(h.Sum(nil))

	return newFileName
}
