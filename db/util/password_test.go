package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	password := RandomString(6)
	// Happy Path
	hashedPassword1, err := HashPassword(password)
	require.NoError(t, err) 
	require.NotEmpty(t, hashedPassword1) // 这里要求 hasedPassword 不能为空

	err = CheckPassword(password, hashedPassword1)
	require.NoError(t, err)

	// 使用错误密码
	wrongPassword := RandomString(6)
	err = CheckPassword(wrongPassword, hashedPassword1)
	require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

	// 验证两个相同的明文密码分别进行哈希，两个哈希值不同
	hashedPassword2, err := HashPassword(password)
	require.NoError(t, err) 
	require.NotEmpty(t, hashedPassword2) 
	require.NotEqual(t, hashedPassword1, hashedPassword2)

}
