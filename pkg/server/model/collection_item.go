package model

import (
	"database/sql"
	"log"
)

// CollectionItem collection_itemテーブルデータ
type CollectionItem struct {
	ID     string
	Name   string
	Rarity int
}

type CollectionItemRepository struct {
	Conn *sql.DB
}

func NewCollectionItemRepository(conn *sql.DB) *CollectionItemRepository {
	return &CollectionItemRepository{
		Conn: conn,
	}
}

type CollectionItemRepositoryInterface interface {
	SelectCollectionItemAll() ([]*CollectionItem, error)
}

var _ CollectionItemRepositoryInterface = (*CollectionItemRepository)(nil)

// SelectCollectionItemAll コレクションアイテムを全取得する
func (r *CollectionItemRepository) SelectCollectionItemAll() ([]*CollectionItem, error) {
	stmt, err := r.Conn.Prepare("SELECT * from collection_item")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	return convertToCollectionItems(rows)
}

// convertToCollectionItems rowsデータをCollectionItemのスライスへ変換する
func convertToCollectionItems(rows *sql.Rows) ([]*CollectionItem, error) {
	defer rows.Close()

	var (
		collectionItems []*CollectionItem
		err             error
	)

	for rows.Next() {
		collectionItem := CollectionItem{}
		if err = rows.Scan(&collectionItem.ID, &collectionItem.Name, &collectionItem.Rarity); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			log.Println(err)
			return nil, err
		}
		collectionItems = append(collectionItems, &collectionItem)
	}
	return collectionItems, err
}
