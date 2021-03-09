package server

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRankingListIntegration(t *testing.T) {
	// モックサーバー
	mux := http.NewServeMux()
	mux.HandleFunc("/test/ranking/list", get(testAuthMiddleware.Authenticate(testRankingHandler.HandleRankingList)))
	server := httptest.NewServer(mux)
	defer server.Close()

	type request struct {
		method  string
		pattern string
		token   string
	}
	type want struct {
		statusCode int
		body       string
	}
	tests := []struct {
		name    string
		before  func()
		request request
		after   func(*http.Response)
		want    want
	}{
		{
			name: "正常:ランキング取得",
			before: func() {
				// シードの作成
				query := `INSERT INTO user(id, auth_token, name, high_score, coin) VALUES ("id1", "token1", "name1", 100, 10000000), ("id2", "token2", "name2", 1000, 10000000), ("id3", "token3", "name3", 100000, 10000000)`
				_, err := testUserRepository.Conn.Exec(query)
				if err != nil {
					t.Errorf("db.TestConn.Exec failed %s", err)
					return
				}
			},
			request: request{
				method:  "GET",
				pattern: "/test/ranking/list?start=1",
				token:   "token1",
			},
			after: func(res *http.Response) {
				// シードの削除
				query := `DELETE FROM user WHERE id in ("id1", "id2", "id3")`
				_, err := testUserRepository.Conn.Exec(query)
				if err != nil {
					t.Errorf("db.TestConn.Exec failed %s", err)
					return
				}
				defer res.Body.Close()
			},
			want: want{
				statusCode: http.StatusOK,
				body: `{
						  "ranks": [
							{
							  "userId": "id3",
							  "userName": "name3",
							  "rank": 1,
							  "score": 100000
							},
							{
							  "userId": "id2",
							  "userName": "name2",
							  "rank": 2,
							  "score": 1000
							},
							{
							  "userId": "id1",
							  "userName": "name1",
							  "rank": 3,
							  "score": 100
							}
						  ]
						}`,
			},
		},
		{
			name: "異常:無効なトークン",
			before: func() {
				// シードの作成
				query := `INSERT INTO user(id, auth_token, name, high_score, coin) VALUES ("id1", "token1", "name1", 100, 10000000), ("id2", "token2", "name2", 1000, 10000000), ("id3", "token3", "name3", 100000, 10000000)`
				_, err := testUserRepository.Conn.Exec(query)
				if err != nil {
					t.Errorf("db.TestConn.Exec failed %s", err)
				}
			},
			request: request{
				method:  "GET",
				pattern: "/test/ranking/list?start=1",
				token:   "invalid-token",
			},
			after: func(res *http.Response) {
				// シードの削除
				query := `DELETE FROM user WHERE id in ("id1", "id2", "id3")`
				_, err := testUserRepository.Conn.Exec(query)
				if err != nil {
					t.Errorf("db.TestConn.Exec failed %s", err)
				}
				defer res.Body.Close()
			},
			want: want{
				statusCode: http.StatusInternalServerError,
				body: `{
							"code": 500,
							"message": "Internal Server Error"
						}`,
			},
		},
		{
			name: "異常:無効なメソッド",
			before: func() {
				// シードの作成
				query := `INSERT INTO user(id, auth_token, name, high_score, coin) VALUES ("id1", "token1", "name1", 100, 10000000), ("id2", "token2", "name2", 1000, 10000000), ("id3", "token3", "name3", 100000, 10000000)`
				_, err := testUserRepository.Conn.Exec(query)
				if err != nil {
					t.Errorf("db.TestConn.Exec failed %s", err)
				}
			},
			request: request{
				method:  "POST",
				pattern: "/test/ranking/list?start=1",
				token:   "token1",
			},
			after: func(res *http.Response) {
				// シードの削除
				query := `DELETE FROM user WHERE id in ("id1", "id2", "id3")`
				_, err := testUserRepository.Conn.Exec(query)
				if err != nil {
					t.Errorf("db.TestConn.Exec failed %s", err)
				}
				defer res.Body.Close()
			},
			want: want{
				statusCode: http.StatusMethodNotAllowed,
				body:       `Method Not Allowed`,
			},
		},
		{
			name: "異常:クエリパラメータエラー",
			before: func() {
				// シードの作成
				query := `INSERT INTO user(id, auth_token, name, high_score, coin) VALUES ("id1", "token1", "name1", 100, 10000000), ("id2", "token2", "name2", 1000, 10000000), ("id3", "token3", "name3", 100000, 10000000)`
				_, err := testUserRepository.Conn.Exec(query)
				if err != nil {
					t.Errorf("db.TestConn.Exec failed %s", err)
				}
			},
			request: request{
				method:  "GET",
				pattern: "/test/ranking/list?notstart=1",
				token:   "token1",
			},
			after: func(res *http.Response) {
				// シードの削除
				query := `DELETE FROM user WHERE id in ("id1", "id2", "id3")`
				_, err := testUserRepository.Conn.Exec(query)
				if err != nil {
					t.Errorf("db.TestConn.Exec failed %s", err)
				}
				defer res.Body.Close()
			},
			want: want{
				statusCode: http.StatusBadRequest,
				body: `{
							"code": 400,
							"message": "Bad Request"
						}`,
			},
		},
		{
			name: "異常:開始順位エラー",
			before: func() {
				// シードの作成
				query := `INSERT INTO user(id, auth_token, name, high_score, coin) VALUES ("id1", "token1", "name1", 100, 10000000), ("id2", "token2", "name2", 1000, 10000000), ("id3", "token3", "name3", 100000, 10000000)`
				_, err := testUserRepository.Conn.Exec(query)
				if err != nil {
					t.Errorf("db.TestConn.Exec failed %s", err)
				}
			},
			request: request{
				method:  "GET",
				pattern: "/test/ranking/list?start=-1",
				token:   "token1",
			},
			after: func(res *http.Response) {
				// シードの削除
				query := `DELETE FROM user WHERE id in ("id1", "id2", "id3")`
				_, err := testUserRepository.Conn.Exec(query)
				if err != nil {
					t.Errorf("db.TestConn.Exec failed %s", err)
				}
				defer res.Body.Close()
			},
			want: want{
				statusCode: http.StatusBadRequest,
				body: `{
							"code": 400,
							"message": "Bad Request"
						}`,
			},
		},
		{
			name: "異常:ユーザ数0",
			before: func() {
			},
			request: request{
				method:  "GET",
				pattern: "/test/ranking/list?start=1",
				token:   "token1",
			},
			after: func(res *http.Response) {
			},
			want: want{
				statusCode: http.StatusInternalServerError,
				body: `{
							"code": 500,
							"message": "Internal Server Error"
						}`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.before() // シード作成

			// リクエスト
			req, err := http.NewRequest(tt.request.method, server.URL+tt.request.pattern, nil)
			if err != nil {
				t.Errorf("http.NewRequest faild: %v", err)
				return
			}
			req.Header.Set("x-token", tt.request.token)

			// 実行してレスポンスを取得
			client := http.DefaultClient
			res, err := client.Do(req)
			if err != nil {
				t.Errorf("http.DefaultClient.Do failed: %v", err)
				return
			}

			if res.StatusCode != tt.want.statusCode {
				t.Errorf("status code = %d, want %d", res.StatusCode, tt.want.statusCode)
			}

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("ioutil.ReadAll failed: %v", err)
				return
			}

			boolean, err := deepEqualString(string(body), tt.want.body)
			if err != nil {
				t.Errorf("response.DeepEqualString() failed %s", err)
			}
			if !boolean {
				t.Errorf("response body = \n%s\n, want \n%s\n", string(body), tt.want.body)
			}

			tt.after(res)
		})
	}
}
