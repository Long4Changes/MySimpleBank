package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	mockdb "github.com/Long4Changes/MySimpleBank/db/mock"
	db "github.com/Long4Changes/MySimpleBank/db/sqlc"
	"github.com/Long4Changes/MySimpleBank/db/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCreateUserAPI(t *testing.T) {
	password := util.RandomString(6)
	hashedPassword, err := util.HashPassword(password)
	require.NoError(t, err)
	user := randomUser(hashedPassword)

	arg := db.CreateUserParams{
		Username:       user.Username,
		HashedPassword: user.HashedPassword,
		FullName:       user.FullName,
		Email:          user.Email,
	}

	testCases := []struct {
		name          string
		username      string
		password      string
		fullName      string
		email         string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, record *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			username: user.Username,
			password: password,
			fullName: user.FullName,
			email:    user.Email,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					// 由于程序执行到 createUser 的时候会在内部再进行一次 hash
					// 由于 random salt 机制，同一明文密码每次进行 hash 得到的 hash 都是不一样的
					// 所以这里传入的 hashedPassword 和底下实际执行程序时传入的不一样
					// 通过 EqCreateUserParams 传入 arg = 期望模版，password = 明文密码
					CreateUser(gomock.Any(), EqCreateUserParams(arg, password)).
					Times(1).
					Return(user, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchUser(t, recorder.Body, user)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			store := mockdb.NewMockStore(ctrl)

			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := "/users"
			body, err := json.Marshal(createUserRequest{
				Username: tc.username,
				Password: tc.password,
				FullName: tc.fullName,
				Email:    tc.email,
			})
			require.NoError(t, err)
			require.NotEmpty(t, body)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}

// 定义 eqCreateUserParamsMatcher 这个类
// 在下面我们让这个类实现了 gomock.Matcher 接口要求的方法
// 只要一个类实现了接口全部方法，Go 就会自动把它当这个接口用
type eqCreateUserParamsMatcher struct {
	arg      db.CreateUserParams
	password string
}

// 这里的输入 x 是真实调用时的实参，gomock 会自动塞进来
func (e eqCreateUserParamsMatcher) Matches(x any) bool {
	arg, ok := x.(db.CreateUserParams)
	if !ok {
		return false
	}

	// 检查 arg 的 HashedPassword 是否是由明文密码 password 经过 hash 得到的
	if err := util.CheckPassword(e.password, arg.HashedPassword); err != nil {
		return false
	}
	// 如果验证确实是，就把期望模版 arg 中的 hash 替换成真实参数里的 hash
	e.arg.HashedPassword = arg.HashedPassword
	// 通过 Deep Equal 深度比较两个结构体是否相等
	// truct values are deeply equal if their corresponding fields, both exported and unexported, are deeply equal.
	return reflect.DeepEqual(e.arg, arg)
}

func (e eqCreateUserParamsMatcher) String() string {
	return fmt.Sprintf("matches arg %v and password %v", e.arg, e.password)
}

// 返回 gomock.Matcher 这个接口类型
func EqCreateUserParams(arg db.CreateUserParams, password string) gomock.Matcher {
	return eqCreateUserParamsMatcher{arg, password}
}

func randomUser(hashedPassword string) db.User {
	return db.User{
		Username:       util.RandomString(6),
		HashedPassword: hashedPassword,
		FullName:       util.RandomString(6),
		Email:          util.RandomEmail(),
	}
}

func requireBodyMatchUser(t *testing.T, body *bytes.Buffer, user db.User) {
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotUser createUserResponse

	err = json.Unmarshal(data, &gotUser)
	require.NoError(t, err)
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.FullName, gotUser.FullName)
	require.Equal(t, user.Email, gotUser.Email)
	require.Equal(t, user.PasswordChangedAt, gotUser.PasswordChangedAt)
	require.Equal(t, user.CreatedAt, gotUser.CreatedAt)
}
