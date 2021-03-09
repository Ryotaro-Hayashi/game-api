package model

import (
	"database/sql"
	"log"
)

type GachaProbability struct {
	CollectionItemId string
	Ratio            int
}

type GachaProbabilityRepository struct {
	Conn *sql.DB
}

func NewGachaRepositoryRepository(conn *sql.DB) *GachaProbabilityRepository {
	return &GachaProbabilityRepository{
		Conn: conn,
	}
}

type GachaProbabilityRepositoryInterface interface {
	SelectGachaProbabilityAll() ([]*GachaProbability, error)
}

var _ GachaProbabilityRepositoryInterface = (*GachaProbabilityRepository)(nil)

// SelectGachaProbabilityAll ガチャ排出確率情報を全取得する
func (r *GachaProbabilityRepository) SelectGachaProbabilityAll() ([]*GachaProbability, error) {
	stmt, err := r.Conn.Prepare("SELECT * FROM gacha_probability")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	return convertToGachaProbabilities(rows)
}

// convertToGachaProbabilities rowsデータをGachaProbabilityのスライスへ変換する
func convertToGachaProbabilities(rows *sql.Rows) ([]*GachaProbability, error) {
	defer rows.Close()

	var (
		gachaProbabilities []*GachaProbability
		err                error
	)

	for rows.Next() {
		gachaProbability := GachaProbability{}
		if err = rows.Scan(&gachaProbability.CollectionItemId, &gachaProbability.Ratio); err != nil {
			if err == sql.ErrNoRows {
				return nil, nil
			}
			log.Println(err)
			return nil, err
		}
		gachaProbabilities = append(gachaProbabilities, &gachaProbability)
	}

	return gachaProbabilities, err
}
