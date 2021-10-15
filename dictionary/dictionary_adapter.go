package dictionary

import (
	"fmt"

	"github.com/awatercolorpen/olap-sql/api/models"
)

// AdapterDataType Adapter数据类型
type AdapterDataType string

const (
	MetricsType AdapterDataType = "metrics"
)

// DataSaveCenter 用于保存指标的逻辑数据信息
type DictionaryAdapter struct {
	// TODO
	set        []*models.DataSet
	sources    []*models.DataSource
	metrics    []*models.Metric
	dimensions []*models.Dimension
}

type DictionaryAdapterOption struct {
	// TODO

}

func NewDictionaryAdapter(option *DictionaryAdapterOption) (*DictionaryAdapter, error) {
	// TODO
	return nil, nil
}

func (d *DictionaryAdapter) GetDataSetByName(name string) (*models.DataSet, error) {
	for _, data := range d.set {
		if data.Name == name {
			return checkDataSetActive(data)
		}
	}
	return nil, fmt.Errorf("can not find '%v' data set", name)
}

func (d *DictionaryAdapter) GetSourcesByIds(ids []uint64) ([]*models.DataSource, error) {
	idsMap := make(map[uint64]bool)
	for _, id := range ids {
		idsMap[id] = true
	}
	// TODO
	// d.db.Preload("Metrics").Preload("Dimensions").Find(&t.sources, "id IN ?", id).Error
	// 满足sources的id在ids，且同时存在于metrics表和dimensions表中。

	return nil, nil
}

func (d *DictionaryAdapter) GetMetricsByIds(ids []uint64) ([]*models.Metric, error) {
	// TODO
	// d.db.Find(&t.metrics, "data_source_id IN ?", id).Error
	idsMap := make(map[uint64]bool)
	for _, id := range ids {
		idsMap[id] = true
	}
	metrics := make([]*models.Metric, 0)
	for _, metric := range d.metrics {
		if _, ok := idsMap[metric.DataSourceID]; ok {
			metrics = append(metrics, metric)
		}
	}
	return metrics, nil
}

func (d *DictionaryAdapter) GetDimensionsByIds(ids []uint64) ([]*models.Dimension, error) {
	// TODO
	// d.db.Find(&t.dimensions, "data_source_id IN ?", id).Error
	idsMap := make(map[uint64]bool)
	for _, id := range ids {
		idsMap[id] = true
	}
	dimensions := make([]*models.Dimension, 0)
	for _, dimension := range d.dimensions {
		if _, ok := idsMap[dimension.DataSourceID]; ok {
			dimensions = append(dimensions, dimension)
		}
	}
	return dimensions, nil
}

func checkDataSetActive(set *models.DataSet) (*models.DataSet, error) {
	if set.Schema == nil {
		return nil, fmt.Errorf("schema is nil for data_set %v", set.Name)
	}
	return set, nil
}

func isValidJoinOns(joinOns models.JoinOns) (id1, id2 uint64, err error) {
	//in1, in2 := joinOns.ID()

	var out1, out2 []uint64

	// TODO get out1, out2
	// db.Table(DefaultOlapSqlModelDimensionTableName).Select("data_source_id").Group("data_source_id").Find(&out1, "id IN ?", in1).Error; err != nil {
	// db.Table(DefaultOlapSqlModelDimensionTableName).Select("data_source_id").Group("data_source_id").Find(&out2, "id IN ?", in2).Error; err != nil {

	if len(out1) != 1 {
		return 0, 0, fmt.Errorf("invalid data_source_id=%v", out1)
	}
	if len(out2) != 1 {
		return 0, 0, fmt.Errorf("invalid data_source_id=%v", out2)
	}
	id1 = out1[0]
	id2 = out2[0]
	return
}

func isValidSecondary(secondary *models.Secondary) error {
	id1, id2, err := isValidJoinOns(models.JoinOns(secondary.JoinOn))
	if err != nil {
		return err
	}
	if id1 != secondary.DataSourceID1 {
		return fmt.Errorf("unmatched data_source_ids, %v != %v", id1, secondary.DataSourceID1)
	}
	if id2 != secondary.DataSourceID2 {
		return fmt.Errorf("unmatched data_source_ids, %v != %v", id2, secondary.DataSourceID2)
	}
	return nil
}

func isValidDataSetSchema(schema *models.DataSetSchema) error {
	if _, err := schema.Tree(); err != nil {
		return err
	}

	for _, v := range schema.Secondary {
		if err := isValidSecondary(v); err != nil {
			return err
		}
	}
	return nil
}
