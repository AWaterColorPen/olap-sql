package models

import (
	"database/sql/driver"
	"fmt"
	"time"

	"gorm.io/gorm"
)

var DefaultOlapSqlModelDataSetTableName = "olap_sql_model_data_sets"

type DataSet struct {
	ID          uint64         `gorm:"column:id;primaryKey"   json:"id,omitempty"`
	CreatedAt   time.Time      `gorm:"column:created_at"      json:"created_at,omitempty"`
	UpdatedAt   time.Time      `gorm:"column:updated_at"      json:"updated_at,omitempty"`
	DeletedAt   gorm.DeletedAt `gorm:"column:delete_at;index" json:"-"`
	Name        string         `gorm:"column:name;unique"     json:"name"`
	Description string         `gorm:"column:description"     json:"description"`
	Schema      *DataSetSchema `gorm:"column:schema"          json:"schema"`
}

func (DataSet) TableName() string {
	return DefaultOlapSqlModelDataSetTableName
}

type DataSetSchema struct {
	PrimaryID uint64       `json:"primary_id"`
	Secondary []*Secondary `json:"secondary"`
}

func (d *DataSetSchema) DataSourceID() []uint64 {
	id := []uint64{d.PrimaryID}
	for _, v := range d.Secondary {
		id = append(id, v.DataSourceID1, v.DataSourceID2)
	}
	return id
}

func (d *DataSetSchema) Tree() (map[uint64][]uint64, error) {
	// 如果没有副数据源，则直接返回主数据源
	if len(d.Secondary) == 0 {
		return map[uint64][]uint64{d.PrimaryID: {}}, nil
	}
	// 校验副数据源关系是否正确，并返回关联关系
	inDegree := map[uint64]int{}
	for _, v := range d.Secondary {
		if _, ok := inDegree[v.DataSourceID1]; !ok {
			inDegree[v.DataSourceID1] = 0
		}
		inDegree[v.DataSourceID2] = inDegree[v.DataSourceID2] + 1
	}

	var root []uint64
	for k, v := range inDegree {
		if v > 1 {
			return nil, fmt.Errorf("it is graph, not a tree. id=%v has indegree=%v", k, v)
		}
		if v == 0 {
			root = append(root, k)
		}
	}
	if len(root) != 1 {
		return nil, fmt.Errorf("there are mulit root. id=%v", root)
	}

	if len(root) == 1 && root[0] != d.PrimaryID {
		return nil, fmt.Errorf("root is not matched %v != %v", root[0], d.PrimaryID)
	}

	tree := map[uint64][]uint64{}
	for _, v := range d.Secondary {
		tree[v.DataSourceID1] = append(tree[v.DataSourceID1], v.DataSourceID2)
	}
	return tree, nil
}

func (d *DataSetSchema) Valid(db *gorm.DB) error {
	if _, err := d.Tree(); err != nil {
		return err
	}

	for _, v := range d.Secondary {
		if err := v.Valid(db); err != nil {
			return err
		}
	}

	return nil
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (d *DataSetSchema) Scan(value interface{}) error {
	return scan(value, d)
}

// Value return json value, implement driver.Valuer interface
func (d DataSetSchema) Value() (driver.Value, error) {
	return value(d)
}

type Secondary struct {
	DataSourceID1 uint64    `json:"data_source_id1"`
	DataSourceID2 uint64    `json:"data_source_id2"`
	JoinOn        []*JoinOn `json:"join_on"`
}

func (s *Secondary) Valid(db *gorm.DB) error {
	id1, id2, err := JoinOns(s.JoinOn).Valid(db)
	if err != nil {
		return err
	}
	if id1 != s.DataSourceID1 {
		return fmt.Errorf("unmatched data_source_ids, %v != %v", id1, s.DataSourceID1)
	}
	if id2 != s.DataSourceID2 {
		return fmt.Errorf("unmatched data_source_ids, %v != %v", id2, s.DataSourceID2)
	}
	return nil
}

type JoinOn struct {
	DimensionID1 uint64 `json:"dimension_id1"`
	DimensionID2 uint64 `json:"dimension_id2"`
}

type JoinOns []*JoinOn

func (j JoinOns) ID() (id1, id2 []uint64) {
	for _, v := range j {
		id1 = append(id1, v.DimensionID1)
		id2 = append(id2, v.DimensionID2)
	}
	return
}

func (j JoinOns) Valid(db *gorm.DB) (id1, id2 uint64, err error) {
	in1, in2 := j.ID()

	var out1, out2 []uint64
	if err = db.Table(DefaultOlapSqlModelDimensionTableName).Select("data_source_id").Group("data_source_id").Find(&out1, "id IN ?", in1).Error; err != nil {
		return
	}
	if len(out1) != 1 {
		return 0, 0, fmt.Errorf("invalid data_source_id=%v", out1)
	}

	if err = db.Table(DefaultOlapSqlModelDimensionTableName).Select("data_source_id").Group("data_source_id").Find(&out2, "id IN ?", in2).Error; err != nil {
		return
	}
	if len(out2) != 1 {
		return 0, 0, fmt.Errorf("invalid data_source_id=%v", out2)
	}
	id1 = out1[0]
	id2 = out2[0]
	return
}
