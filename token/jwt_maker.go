package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const minSecretKeySize = 32

// JWTMaker is a JSON Web Token maker
// JWTMaker 必须要实现 Maker 接口
type JWTMaker struct {
	secretKey string
}

func NewJWTMaker(secretKey string) (Maker, error) {
	// 设置一个最短密钥长度
	if len(secretKey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid key size: must be at least %d characters", minSecretKeySize)
	}

	// 这里和 custom gomock matcher 那里的写法差不多
	// 想要直接返回这个接口类型，JWTMaker 就要实现 Maker 接口的所有方法
	return &JWTMaker{secretKey}, nil
}

// CreateToken creates a new token for a specific username and duration
func (maker *JWTMaker) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	// Claims 是载荷字段，也就是 token 里携带的数据，Payload 就是 claims 的结构体
	// *Payload does not implement jwt.Claims (missing method Valid)
	// 这里想要传入 payload，那么 Payload 类型也要实现 jwt.Claims 接口的方法，也就是 Valid()
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	// 使用 secretKey 去给 jwtToken 签名
	return jwtToken.SignedString([]byte(maker.secretKey))
}

// VerifyToken checks if the token is valid or not
func (maker *JWTMaker) VerifyToken(token string) (*Payload, error) {
	// keyFunc 的作用是去验证请求头中记录的 “alg” 是否和 Signing 时所用到的算法一致，防止算法混淆攻击
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		// 由于上面我们用的算法是 HS256，var jwt.SigningMethodHS256 *jwt.SigningMethodHMAC
		// 这里要确认要验证的 token 是否是 HMAC 家族的
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrorInvalidToken
		}

		// 如果匹配上了，就返回密钥
		return []byte(maker.secretKey), nil
	}
	// 
	jwtToken, err := jwt.ParseWithClaims(token, &Payload{}, keyFunc)
	if err != nil {
		// 将传回来的 err 解析成可细分原因的 JWT 错误
		verr, ok := err.(*jwt.ValidationError)
		if ok && errors.Is(verr.Inner, ErrorExpireToken) {
			return nil, ErrorExpireToken
		}
		return nil, ErrorInvalidToken
	}

	payload, ok := jwtToken.Claims.(*Payload)
	if !ok {
		return nil, ErrorInvalidToken
	} 
	return payload, nil
}
