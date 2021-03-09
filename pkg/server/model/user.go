//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=./mock_$GOPACKAGE/mock_$GOFILE

package model

import (
	"database/sql"
	"log"
)

// User userテーブルデータ
type User struct {
	ID        string
	AuthToken string
	Name      string
	HighScore int
	Coin      int
}

type UserRepository struct {
	Conn *sql.DB
}

func NewUserRepository(conn *sql.DB) *UserRepository {
	return &UserRepository{
		Conn: conn,
	}
}

type UserRepositoryInterface interface {
	InsertUser(record *User) error
	SelectUserByAuthToken(authToken string) (*User, error)
	UpdateUserByPrimaryKey(record *User) error
	SelectUserByPrimaryKey(userID string) (*User, error)
	UpdateUserCoinAndHighScoreByPrimaryKey(id string, coin int, highScore int) error
	SelectUsersOrderByHighScoreDesc(limit int, offset int) ([]*User, error)
	UpdateUserCoinByPrimaryKey(tx *sql.Tx, userID string, coin int) error
	SelectUserByPrimaryKeyForUpdate(tx *sql.Tx, userID string) (*User, error)
}

// インターフェースを満たしているかを確認
var _ UserRepositoryInterface = (*UserRepository)(nil)

// InsertUser データベースをレコードを登録する
func (r *UserRepository) InsertUser(record *User) error {
	stmt, err := r.Conn.Prepare("INSERT INTO user(id, auth_token, name, high_score, coin) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(record.ID, record.AuthToken, record.Name, record.HighScore, record.Coin)
	return err
}

// SelectUserByAuthToken auth_tokenを条件にレコードを取得する
func (r *UserRepository) SelectUserByAuthToken(authToken string) (*User, error) {
	row := r.Conn.QueryRow("SELECT * from user WHERE auth_token = ?", authToken)
	return convertToUser(row)
}

// SelectUserByPrimaryKey 主キーを条件にレコードを取得する
func (r *UserRepository) SelectUserByPrimaryKey(userID string) (*User, error) {
	row := r.Conn.QueryRow("SELECT * from user WHERE id = ?", userID)
	return convertToUser(row)
}

// UpdateUserByPrimaryKey 主キーを条件にレコードを更新する
func (r *UserRepository) UpdateUserByPrimaryKey(record *User) error {
	stmt, err := r.Conn.Prepare("UPDATE user SET name = ? where  id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(record.Name, record.ID)
	return err
}

// UpdateUserCoinAndHighScoreByPrimaryKey 主キーを条件に所持コインとハイスコアを更新する
func (r *UserRepository) UpdateUserCoinAndHighScoreByPrimaryKey(id string, coin int, highScore int) error {
	stmt, err := r.Conn.Prepare("Update user SET coin = ?, high_score = ? where id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(coin, highScore, id)

	return err
}

// SelectUsersOrderByHighScoreDesc ハイスコア順に指定順位から指定件数を取得する
func (r *UserRepository) SelectUsersOrderByHighScoreDesc(limit int, offset int) ([]*User, error) {
	stmt, err := r.Conn.Prepare("SELECT * FROM user ORDER BY high_score DESC LIMIT ? OFFSET ?")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(limit, offset-1)
	if err != nil {
		return nil, err
	}

	return convertToUsers(rows)
}

// UpdateUserCoinByPrimaryKey 主キーを条件にコインを更新する
func (r *UserRepository) UpdateUserCoinByPrimaryKey(tx *sql.Tx, userID string, coin int) error {
	stmt, err := tx.Prepare("UPDATE user SET coin = ? where  id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(coin, userID)
	return err
}

// SelectUserByPrimaryKeyForUpdate 主キーを条件に排他ロックでユーザ情報を取得する
func (r *UserRepository) SelectUserByPrimaryKeyForUpdate(tx *sql.Tx, userID string) (*User, error) {
	row := tx.QueryRow("SELECT * from user WHERE id = ? FOR UPDATE", userID)
	return convertToUser(row)
}

// convertToUser rowデータをUserデータへ変換する
func convertToUser(row *sql.Row) (*User, error) {
	user := User{}
	err := row.Scan(&user.ID, &user.AuthToken, &user.Name, &user.HighScore, &user.Coin)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Println(err)
		return nil, err
	}
	return &user, nil
}

// convertToUsers rowsデータをUserのスライスへ変換する
func convertToUsers(rows *sql.Rows) ([]*User, error) {
	defer rows.Close()

	var (
		users []*User
		err   error
	)

	for rows.Next() {
		user := User{}
		if err = rows.Scan(&user.ID, &user.AuthToken, &user.Name, &user.HighScore, &user.Coin); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			log.Println(err)
			return nil, err
		}
		users = append(users, &user)
	}
	return users, err
}
