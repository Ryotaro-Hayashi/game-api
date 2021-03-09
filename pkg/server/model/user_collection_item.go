package model

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
)

// UserCollectionItem user_collection_itemテーブルデータ
type UserCollectionItem struct {
	UserID           string
	CollectionItemID string
}

type UserCollectionItemRepository struct {
	Conn *sql.DB
}

func NewUserCollectionItemRepository(conn *sql.DB) *UserCollectionItemRepository {
	return &UserCollectionItemRepository{
		Conn: conn,
	}
}

type UserCollectionItemRepositoryInterface interface {
	SelectUserCollectionItemsByUserID(userID string) ([]*UserCollectionItem, error)
	BulkInsertUserCollectionItem(tx *sql.Tx, newCollectionItemSlice []*UserCollectionItem) error
}

var _ UserCollectionItemRepositoryInterface = (*UserCollectionItemRepository)(nil)

// SelectUserCollectionItemsByUserID ユーザIDを条件に所持アイテムを取得する
func (r *UserCollectionItemRepository) SelectUserCollectionItemsByUserID(userID string) ([]*UserCollectionItem, error) {
	stmt, err := r.Conn.Prepare("SELECT * FROM user_collection_item WHERE user_id = ?")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query(userID)
	if err != nil {
		return nil, err
	}

	return convertToUserCollectionItems(rows)
}

// BulkInsertUserCollectionItem Newアイテムを登録する
func (r *UserCollectionItemRepository) BulkInsertUserCollectionItem(tx *sql.Tx, newCollectionItemSlice []*UserCollectionItem) error {

	placeholder := make([]string, 0, len(newCollectionItemSlice))
	queryArgs := make([]interface{}, 0, len(newCollectionItemSlice)*2)
	for _, newCollectionItem := range newCollectionItemSlice {
		placeholder = append(placeholder, "(?, ?)")
		queryArgs = append(queryArgs, newCollectionItem.UserID, &newCollectionItem.CollectionItemID)
	}

	query := fmt.Sprintf("INSERT INTO user_collection_item (user_id, collection_item_id) VALUES %s", strings.Join(placeholder, ", "))
	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(queryArgs...)
	if err != nil {
		return err
	}

	return err
}

// convertToUserCollectionItems rowsデータをUserCollectionItemのスライスへ変換する
func convertToUserCollectionItems(rows *sql.Rows) ([]*UserCollectionItem, error) {
	defer rows.Close()

	var (
		userCollectionItems []*UserCollectionItem
		err                 error
	)

	for rows.Next() {
		userCollectionItem := UserCollectionItem{}
		if err = rows.Scan(&userCollectionItem.UserID, &userCollectionItem.CollectionItemID); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			log.Println(err)
			return nil, err
		}
		userCollectionItems = append(userCollectionItems, &userCollectionItem)
	}
	return userCollectionItems, err
}
