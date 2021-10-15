package dictionary

import (
	"fmt"

	"github.com/awatercolorpen/olap-sql/api/models"
)

type AdapterType string

const (
	dbAdapter   AdapterType = "DB"
	fileAdapter AdapterType = "FILE"
)

// Adapter Adapter适配器
type Adapter interface {
	// TODO 有待考虑
	GetDataSetByName(string) (interface{}, error)
	GetSourcesByIds([]interface{}) ([]interface{}, error)
	GetMetricsByIds([]interface{}) ([]interface{}, error)
	GetDimensionsByIds([]interface{}) ([]interface{}, error)
}

// AdapterOption Adapter配置
type AdapterOption struct {
	adapterType AdapterType
	dsn         string
}

func NewAdapter(option *AdapterOption) (Adapter, error) {
	// 根据不同的Type去实例化不同的Adapter
	switch option.adapterType {
	case dbAdapter:
		// TODO
	case fileAdapter:
		// TODO
	}
	return nil, nil
}

// DataSaveCenter 用于保存指标的逻辑数据信息
type DictionaryAdapter struct {
	// TODO
	set        []*models.DataSet
	sources    []*models.DataSource
	metrics    []*models.Metric
	dimensions []*models.Dimension
}

func NewDictionaryAdapter(option *AdapterOption) (*DictionaryAdapter, error) {
	// TODO
	// 1. 根据Option解析文件读入数据
	// 2. CheckValid调用
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
	idsMap := getIdsMap(ids)

	metricsSourcesIdsMap := make(map[uint64]bool)
	for _, metric := range d.metrics {
		metricsSourcesIdsMap[metric.DataSourceID] = true
	}

	dimensionsSourcesIdsMap := make(map[uint64]bool)
	for _, dimension := range d.dimensions {
		dimensionsSourcesIdsMap[dimension.DataSourceID] = true
	}

	result := make([]*models.DataSource, 0)
	for _, source := range d.sources {
		_, ok := idsMap[source.ID]
		_, ok2 := metricsSourcesIdsMap[source.ID]
		_, ok3 := dimensionsSourcesIdsMap[source.ID]
		if ok && ok2 && ok3 {
			result = append(result, source)
		}
	}
	return result, nil
}

func (d *DictionaryAdapter) GetMetricsByIds(ids []uint64) ([]*models.Metric, error) {
	idsMap := getIdsMap(ids)
	metrics := make([]*models.Metric, 0)
	for _, metric := range d.metrics {
		if _, ok := idsMap[metric.DataSourceID]; ok {
			metrics = append(metrics, metric)
		}
	}
	return metrics, nil
}

func (d *DictionaryAdapter) GetDimensionsByIds(ids []uint64) ([]*models.Dimension, error) {
	idsMap := getIdsMap(ids)
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

func (d *DictionaryAdapter) isValidJoinOns(joinOns models.JoinOns) (id1, id2 uint64, err error) {
	in1, in2 := joinOns.ID()

	in1Map := getIdsMap(in1)
	in2Map := getIdsMap(in2)

	var out1, out2 map[uint64]bool

	for _, dimension := range d.dimensions {
		id := dimension.ID
		if _, ok := in1Map[id]; ok {
			out1[dimension.DataSourceID] = true
		}
		if _, ok := in2Map[id]; ok {
			out2[dimension.DataSourceID] = true
		}
	}

	if len(out1) != 1 {
		return 0, 0, fmt.Errorf("invalid data_source_id=%v", out1)
	}
	if len(out2) != 1 {
		return 0, 0, fmt.Errorf("invalid data_source_id=%v", out2)
	}

	for id := range out1 {
		id1 = id
	}

	for id := range out2 {
		id2 = id
	}
	return
}

func (d *DictionaryAdapter) isValidSecondary(secondary *models.Secondary) error {
	id1, id2, err := d.isValidJoinOns(models.JoinOns(secondary.JoinOn))
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

func (d *DictionaryAdapter) isValidDataSetSchema(schema *models.DataSetSchema) error {
	if _, err := schema.Tree(); err != nil {
		return err
	}

	for _, v := range schema.Secondary {
		if err := d.isValidSecondary(v); err != nil {
			return err
		}
	}
	return nil
}

func (d *DictionaryAdapter) isValidDataSet(set *models.DataSet) error {
	return d.isValidDataSetSchema(set.Schema)
}

func (d *DictionaryAdapter) isValidAdapterCheck(adapter *Adapter) error {
	for _, set := range d.set {
		if err := d.isValidDataSet(set); err != nil {
			return err
		}
	}
	return nil
}

func getIdsMap(ids []uint64) map[interface{}]interface{} {
	idsMap := make(map[interface{}]interface{})
	for _, id := range ids {
		idsMap[id] = true
	}
	return idsMap
}
