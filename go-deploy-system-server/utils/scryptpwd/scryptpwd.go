package scryptpwd

import (
	"encoding/base64"
	"fmt"
	"go-deploy-system-server/utils"

	"golang.org/x/crypto/scrypt"
)

// 密码加密模块

const (
	PWSALTBYTES = 8  // 加盐长度
	PWHASHBYTES = 32 // 加密长度
)

var PWSALT = utils.PwdKey // 加盐字符串

func ScryptPwd(password string) string {
	salt := make([]byte, PWSALTBYTES)
	salt = []byte(PWSALT)

	// func Key(password, salt []byte, N, r, p, keyLen int) ([]byte, error)
	// 建议参数为 N=32768 r=8 p=1
	// 参数N、r和p应该随着内存延迟和CPU并行度的增加而增加
	hashPw, err := scrypt.Key([]byte(password), salt, 32768, 8, 1, PWHASHBYTES)
	if err != nil {
		fmt.Println("加密失败", err)
		return ""
	}

	return base64.StdEncoding.EncodeToString(hashPw)
}
