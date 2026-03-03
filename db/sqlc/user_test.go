package db

import (
	"context"
	"testing"
	"time"

	"github.com/Long4Changes/MySimpleBank/db/util"
	"github.com/stretchr/testify/require"
)

// 这个测试文件可以在 account_test.go 的基础上改一改
func createRandomUser(t *testing.T) User {
	hashedPassword, err := util.HashPassword(util.RandomString(6))
	require.NoError(t, err)

	arg := CreateUserParams{
		Username:    util.RandomOwner(),
		// "secret" 正确做法是随机生成一个密码，并用 Bcrypt 进行 Hash
		HashedPassword:  hashedPassword,
		FullName: util.RandomOwner(),
		Email: util.RandomEmail(),
	}

	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)      // 检查是否没有产生错误
	require.NotEmpty(t, user) // 检查用户是否为空

	// 检查期待值和实际值是否相等
	require.Equal(t, arg.Username, user.Username)
	require.Equal(t, arg.HashedPassword, user.HashedPassword)
	require.Equal(t, arg.FullName, user.FullName)
	require.Equal(t, arg.Email, user.Email)

	// 检查自动生成的 PasswordChangedAt 是否为零值以及  CreatedAt 是否为空
	require.True(t, user.PasswordChangedAt.IsZero())
	require.NotZero(t, user.CreatedAt)

	return user 
}

func TestCreateUser(t *testing.T) {
	createRandomUser(t)
}

func TestGetUser(t *testing.T) {
	user1 := createRandomUser(t)
	user2, err := testQueries.GetUser(context.Background(), user1.Username)
	require.NoError(t, err)
	require.NotEmpty(t, user2)

	require.Equal(t, user1.Username, user2.Username)
	require.Equal(t, user1.HashedPassword, user2.HashedPassword)
	require.Equal(t, user1.FullName, user2.FullName)
	require.Equal(t, user1.Email, user2.Email)
	require.WithinDuration(t, user1.PasswordChangedAt, user2.PasswordChangedAt, time.Second)
	require.WithinDuration(t, user1.CreatedAt, user2.CreatedAt, time.Second)
}
