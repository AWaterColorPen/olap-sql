package dictionary

import (
	"github.com/awatercolorpen/olap-sql/api/types"
)

type Think struct {
	/* TODO

	TODO
	1. YAML文件的框架格式定义[简单的想法是复现DB中存储的内容]: DB里的内容是什么?
	2. YAML文件的解析方式[已经有方法实现?]
	3. dictionary, translator 逻辑拆分 -> 增加新结构 dataSaveCenter存储解析数据


	DataSaveCenter[数据存放] (目前根据translator中的数据决定)
		set        *models.DataSet
		sources    []*models.DataSource
		metrics    []*models.Metric
		dimensions []*models.Dimension
	需要实现的接口:
		1. 初始化(New DataSaveCenter)
		2. 返回数据(根据query的条件进行return) [原来是直接在数据库里查询，现在是基于set,sources, metrics, dimensions 需要根据query的条件进行过滤信息返回给translator]

	Dictionary_translator
	原本的translator
	*对于一个query(查询请求)，把他转换成一个olap的sql语句，查询数据库里的数据
	【做的是一个转换的任务，然后这个转换的任务需要借助的主要是JoinTree和MetricsMap，本身不需要去记录那些数据]
		1. 初始化(New Translator)
		2. 查询接口用于dictionary进行调用 -> 根据传过来的query然后构建sql语句

	Dictionary
		字典需要存储什么?需要有什么功能?
		*本身应该有字(数据) -> 数据应该从哪里获取(不是字典所需要关心的事情,字典只要有数据就可以了) -> DataSaveCenter[数据存储]
		*提供查询功能，比如查A[ 能够根据A的一些特征去快速查询到A的相关信息 ] -> 具体怎么查询【查询的细节】不是字典所需要关心的事情(只要有这个接口，我调用就拿到数据就行) -> DataFinder[数据查询]
		【功能上来说，拿到一本字典，就是用来查数据】

		Dictionary 包含DataSaveCenter 和 Dictionary_translator两个内容:
			接口: 1. 初始化 -> 调用DataSaveCenter的初始化 -> 提供DataSaveCenter给Dictionary_translator 以进行初始化
				  2. 查询接口 (输入: query{}) -> return result(sql语句)


	*/
}

type DictionaryOption struct {
	AdapterOption
}
type Dictionary struct {
	adapter *DictionaryAdapter
}

func (d *Dictionary) Create(item interface{}) error {
	return d.adapter.Create(item)
}

func (d *Dictionary) Translator(query *types.Query) (Translator, error) {
	d.adapter.fillSourceMetricsAndDimensions()
	set, err := d.adapter.GetDataSetByName(query.DataSetName)
	if err != nil {
		return nil, err
	}

	id := set.Schema.DataSourceID()
	sources, err := d.adapter.GetSourcesByIds(id)
	if err != nil {
		return nil, err
	}

	metrics, err := d.adapter.GetMetricsByIds(id)
	if err != nil {
		return nil, err
	}

	dimensions, err := d.adapter.GetDimensionsByIds(id)
	if err != nil {
		return nil, err
	}

	t := &DictionaryTranslator{
		set:        set,
		sources:    sources,
		metrics:    metrics,
		dimensions: dimensions,
	}
	return t, nil
}

func (d *Dictionary) Translate(query *types.Query) (*types.Request, error) {
	translator, err := d.Translator(query)
	if err != nil {
		return nil, err
	}
	return translator.Translate(query)
}

func NewDictionary(option *DictionaryOption) (*Dictionary, error) {
	// 初始化DictionaryAdapter
	adapter, err := NewDictionaryAdapter(&option.AdapterOption)
	if err != nil {
		return nil, err
	}

	return &Dictionary{adapter: adapter}, nil
}
