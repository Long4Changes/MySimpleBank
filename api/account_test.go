package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	mockdb "github.com/Long4Changes/MySimpleBank/db/mock"
	db "github.com/Long4Changes/MySimpleBank/db/sqlc"
	"github.com/Long4Changes/MySimpleBank/db/util"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestGetAccountAPI(t *testing.T) {
	/*
		// 以下只进行了 Happy Path 完美路径 的测试
		account := randomAccount()

		ctrl := gomock.NewController(t)

		// Go 1.14 之后不再需要手动调用 ctrl.Finish()
		// 现在把 t 传入 gomock.NewController() 的时候，会自动执行类似 t.Cleanup(ctrl.Finish) 的操作
		// defer ctrl.Finish()

		store := mockdb.NewMockStore(ctrl)
		// build stubs
		// 这里 Mock 的方法 GetAccount 的输入输出结构还是要保持和真正的一样，GetAccount(ctx context.Context, id int64) (Account, error)
		store.EXPECT().
			// 希望调用 GetAccount 时具有 任何 上下文和此特定账户的 ID 参数
			// Eq stands for Equal
			GetAccount(gomock.Any(), gomock.Eq(account.ID)).
			Times(1).
			Return(account, nil)

		// start test server and send request
		// 不启动真实的 Web 服务器，而是再内存里“假装”发起了一次 HTTP 请求
		server := NewServer(store)
		// 来自 Go 标准库 net/http/httptest，它实现了 http.ResponseWriter 接口
		// 但它不会把数据写到网络网卡里，而是直接写进自己内部的恶一个内存缓冲区（Buffer）里
		recorder := httptest.NewRecorder()

		url := fmt.Sprintf("/accounts/%d", account.ID)
		// 在内存捏造出一个 HTTP 请求对象，它不需要经过浏览器的封装和操作系统的网络层
		request, err := http.NewRequest(http.MethodGet, url, nil)
		require.NoError(t, err)

		// Gin 的核心引擎 Engine 实现了标准库的 http.Handler 接口
		// 它必须包含一个叫 ServeHTTP 的方法
		// 当真实服务器运行时，是Go 底层的网络模块在收到 TCP 报文后，去调用这个 ServeHTTP
		// 但在测试里，我们人为地、直接地调用了这个方法
		// 调用 ServerHTTP 后，就会去运行 account.go
		server.router.ServeHTTP(recorder, request)
		// check response
		// recorder.Code 存储 response status code 响应状态吗
		// recorder.Body 存储 response body 响应正文
		require.Equal(t, http.StatusOK, recorder.Code)
		requireBodyMatchAccount(t, recorder.Body, account)
	*/

	account := randomAccount()

	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:      "NotFound",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					// 找不到数据，用 ErrNoRows
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					// ErrConnDone 数据库连接已经用完或失效
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidID",
			// 前端传了一个非法的ID，ID的最小值规定为1
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					// 第二个参数可以是任何 ID，即表明无论参数是什么，都不该调用 GetAccount()
					GetAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	// 遍历测试用例表
	for i := range testCases {
		// 用 tc 变量去存储当前测试用例的数据
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			store := mockdb.NewMockStore(ctrl)
			// build stubs
			// 调用 buildStubs 函数
			tc.buildStubs(store)

			// start test server and send request
			server := NewServer(store)
			recorder := httptest.NewRecorder()

			// 这里原来是 account.ID
			// 当测试 InvalidID 用例的时候，这里如果不修改无法通过测试
			// 因为 account 是随机生成的，account.ID 总是一个有效的值
			// 所以这里 account.ID --> tc.accountID
			url := fmt.Sprintf("/accounts/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			// check response
			// 调用 checkResponse 函数
			tc.checkResponse(t, recorder)
		})
	}

}

func TestCreateAccountAPI(t *testing.T) {
	/*
		// 以下只进行了 Happy Path 完美路径 的测试
		owner := util.RandomOwner()
		currency := util.RandomCurrency()

		arg := db.CreateAccountParams{
			Owner:    owner,
			Currency: currency,
			Balance:  0,
		}

		account := db.Account{
			Owner:    owner,
			Currency: currency,
			Balance:  0,
		}

		ctrl := gomock.NewController(t)

		store := mockdb.NewMockStore(ctrl)

		// build stubs
		store.EXPECT().
			CreateAccount(gomock.Any(), gomock.Eq(arg)).
			Times(1).
			Return(account, nil)

		server := NewServer(store)
		recorder := httptest.NewRecorder()

		url := "/accounts"

		// 先把数据 json.Marshal(...) 成 []byte
		body, err := json.Marshal(createAccountRequest{
			Owner:    owner,
			Currency: currency,
		})
		require.NoError(t, err)

		// func http.NewRequest(method string, url string, body io.Reader) (*http.Request, error)
		// 第三个参数要求是 io.Reader 类型
		// 用 bytes.NewReader(body) 当作请求体
		request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
		require.NoError(t, err)
		// request.Header.Set("Content-Type", "application/json")

		server.router.ServeHTTP(recorder, request)

		// check response
		require.Equal(t, http.StatusOK, recorder.Code)
		requireBodyMatchAccount(t, recorder.Body, account)
	*/
	owner := util.RandomOwner()
	currency := util.RandomCurrency()

	arg := db.CreateAccountParams{
		Owner:    owner,
		Currency: currency,
		Balance:  0,
	}

	account := db.Account{
		Owner:    owner,
		Currency: currency,
		Balance:  0,
	}

	testCases := []struct {
		name          string
		Owner         string
		Currency      string
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			Owner:    owner,
			Currency: currency,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		},
		{
			name:     "InternalError",
			Owner:    owner,
			Currency: currency,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:     "InvalidInput",
			Owner:    "",
			Currency: currency,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					CreateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

			url := "/accounts"

			body, err := json.Marshal(createAccountRequest{
				Owner:    tc.Owner,
				Currency: tc.Currency,
			})

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)
			tc.checkResponse(t, recorder)
		})
	}

}

func TestListAccountAPI(t *testing.T) {
	allAccounts := make([]db.Account, util.RandomInt(1, 100))

	for i := range allAccounts {
		account := randomAccount()
		allAccounts[i].ID = account.ID
		allAccounts[i].Owner = account.Owner
		allAccounts[i].Balance = account.Balance
		allAccounts[i].Currency = account.Currency
		allAccounts[i].CreatedAt = account.CreatedAt
	}
	page_size := util.RandomInt32(5, 10)
	totalPages := (int32(len(allAccounts) + int(page_size) - 1)) / page_size
	page_id := util.RandomInt32(1, totalPages)

	arg := db.ListAccountsParams{
		Limit:  page_size,
		Offset: (page_id - 1) * page_size,
	}

	start := int((page_id - 1) * page_size)
	end := start + int(page_size)
	if end > len(allAccounts) {
		end = len(allAccounts)
	}
	resultAccounts := allAccounts[start:end]

	testCases := []struct {
		name          string
		pageID        int32
		pageSize      int32
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:     "OK",
			pageID:   page_id,
			pageSize: page_size,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(resultAccounts, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccounts(t, recorder.Body, resultAccounts)
			},
		},
		{
			name:     "InternalError",
			pageID:   page_id,
			pageSize: page_size,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return([]db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:     "Invalid Input",
			pageID:   page_id,
			pageSize: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					ListAccounts(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

			url := fmt.Sprintf("/accounts?page_id=%d&page_size=%d", tc.pageID, tc.pageSize)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}

func TestUpdateAccountAPI(t *testing.T) {
	account := randomAccount()
	resultBalance := util.RandomMoney() 
	account.Balance = resultBalance

	arg := db.UpdateAccountParams{
		ID: account.ID,
		Balance: resultBalance,
	}
	testCases := []struct{
		name string 
		accountID int64
		accountBalance int64
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK", 
			accountID: account.ID,
			accountBalance: resultBalance,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyMatchAccount(t, recorder.Body, account)
			},
		}, 
		{
			name: "NotFound", 
			accountID: account.ID,
			accountBalance: resultBalance,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		}, 
		{
			name: "InternalError", 
			accountID: account.ID,
			accountBalance: resultBalance,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateAccount(gomock.Any(), gomock.Eq(arg)).
					Times(1).
					Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidInput", 
			accountID: 0,
			accountBalance: resultBalance,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					UpdateAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {

			ctrl := gomock.NewController(t)
			store := mockdb.NewMockStore(ctrl)

			// build stubs
			tc.buildStubs(store)

			server := NewServer(store)
			recorder := httptest.NewRecorder()

			url := "/accounts/update"

			body, err := json.Marshal(UpdateAccountRequest{
				ID: tc.accountID,
				Balance: tc.accountBalance,
			})
			require.NoError(t, err)

			request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			// check response
			tc.checkResponse(t, recorder)

		})
	}
}

func TestDeleteAccountAPI(t *testing.T) {
	account := randomAccount()
	testCases := []struct{
		name string 
		accountID int64
		buildStubs func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "OK", 
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
			},
		}, 
		{
			name: "NotFound", 
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		}, 
		{
			name: "InternalError", 
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteAccount(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name: "InvalidInput", 
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					DeleteAccount(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
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

			url := fmt.Sprintf("/accounts/delete/%d", tc.accountID)
			request, err := http.NewRequest(http.MethodPost, url, nil)
			require.NoError(t, err)

			server.router.ServeHTTP(recorder, request)

			tc.checkResponse(t, recorder)
		})
	}
}
func randomAccount() db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 1000),
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}

func requireBodyMatchAccount(t *testing.T, body *bytes.Buffer, account db.Account) {
	// data, err := ioutil.ReadAll(body)
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccount db.Account
	// 将从响应正文数据获得的账户对象反解到 gotAccount 中
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, account, gotAccount)
}

func requireBodyMatchAccounts(t *testing.T, body *bytes.Buffer, accounts []db.Account) {
	// data, err := ioutil.ReadAll(body)
	data, err := io.ReadAll(body)
	require.NoError(t, err)

	var gotAccounts []db.Account
	// 将从响应正文数据获得的账户对象反解到 gotAccount 中
	err = json.Unmarshal(data, &gotAccounts)
	require.NoError(t, err)
	require.Len(t, gotAccounts, len(accounts))
	for i := range gotAccounts {
		require.Equal(t, accounts[i], gotAccounts[i])
	}
}
