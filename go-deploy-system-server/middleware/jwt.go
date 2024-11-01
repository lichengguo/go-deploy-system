package middleware

import (
	"go-deploy-system-server/utils"
	"go-deploy-system-server/utils/errmsg"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// 中间件

// jwtkey 加盐字符串
var JwtKey = []byte(utils.JwtKey)

// MyCustomClaims token结构体(官方照抄即可)
type MyCustomClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

// SetToken 生成token字符串
func SetToken(name string) (string, int64) {
	expireTime := time.Now().Add(8760 * time.Hour) // token过期时间
	claims := MyCustomClaims{
		name,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(), // 过期时间
			Issuer:    "yiihuaAdmin",     // 签发人
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(JwtKey)
	if err != nil {
		return "", errmsg.ERROR
	}
	return ss, errmsg.SUCCESS
}

// CheckToken 验证token
func CheckToken(tokenStr string) (*MyCustomClaims, int64) {
	token, _ := jwt.ParseWithClaims(tokenStr, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JwtKey, nil
	})

	if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
		return claims, errmsg.SUCCESS
	} else {
		return nil, errmsg.ERROR_TOKEN_FAIL // token验证不通过
	}
}

// JwtToken Jwt登录钩子
func JwtToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int64
		tokenHeader := c.Request.Header.Get("Authorization")

		// token不存在
		if tokenHeader == "" {
			code = errmsg.ERROR_TOKEN_NOT_EXIST
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": errmsg.GetErrMsg(code),
			})
			c.Abort() // 不调用后续的函数
			return
		}

		// token格式错误
		chkToken := strings.Split(tokenHeader, " ")
		if len(chkToken) != 2 || chkToken[0] != "Bearer" {
			code = errmsg.ERROR_TOKEN_TYPE_WRONG
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}

		key, tCode := CheckToken(chkToken[1])

		// token解析失败
		if tCode != errmsg.SUCCESS {
			code = tCode
			c.JSON(http.StatusOK, gin.H{
				"code":    code,
				"message": errmsg.GetErrMsg(code),
			})
			c.Abort()
			return
		}

		// token验证通过
		c.Set("name", key.Name) // 保存用户名到上下文
		c.Next()                // 调用后续函数
	}
}
