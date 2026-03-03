package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns the bcrypt hash of the password
func HashPassword(password string) (string, error) {
	// byte = uint8 一个字节 8 bit
	// rune = int32 一个 Unicode 码点
	// 用户密码明文输入
    // 生成随机 salt（每个密码都不同），也就是说即使两个相同的明文密码去进行哈希，得到的哈希值仍然不同
	// unencodedSalt := make([]byte, maxSaltSize)
	// _, err = io.ReadFull(rand.Reader, unencodedSalt)

	// 把密码和 salt 一起做 bcrypt 计算，迭代次数由 cost 控制
	// 得到最终 hash 字符串，存数据库
	// cost 工作因子（难度系数）用于控制迭代次数（哈希有多慢） 默认值 DefaultCost = 10
	// cost 低：快，但对抗暴力破解的能力弱；cost 高：慢，但更安全
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// %w 对应的参数类型就是 error 类型，作用是包装原始错误，后续可用 errors.Is() / errors.As() / errors.Unwrap() 继续追踪底层错误。
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// CheckPassword checks if the provided password is correct or not 
func CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}