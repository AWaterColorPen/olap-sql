package models

import (
	"fmt"
	"github.com/awatercolorpen/olap-sql/api/types"
)

type DataSet struct {
	Name        string         `toml:"name"             json:"name"`
	DBType      types.DBType   `toml:"type"             json:"type"`
	Description string         `toml:"description"      json:"description"`
	Schema      *DataSetSchema `toml:"schema"           json:"schema"`
}

func (d *DataSet) DataSource() []string {
	var source []string
	return source
}

type DataSetSchema struct {
	PrimaryID uint64       `toml:"primary_id" json:"primary_id"`
	Secondary []*Secondary `toml:"secondary" json:"secondary"`
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

type Secondary struct {
	DataSourceID1 uint64    `toml:"data_source_id1" json:"data_source_id1"`
	DataSourceID2 uint64    `toml:"data_source_id2" json:"data_source_id2"`
	JoinOn        []*JoinOn `toml:"join_on"         json:"join_on"`
}

type JoinOn struct {
	DimensionID1 uint64 `toml:"dimension_id1" json:"dimension_id1"`
	DimensionID2 uint64 `toml:"dimension_id2" json:"dimension_id2"`
}

type JoinOns []*JoinOn

func (j JoinOns) ID() (id1, id2 []uint64) {
	for _, v := range j {
		id1 = append(id1, v.DimensionID1)
		id2 = append(id2, v.DimensionID2)
	}
	return
}
