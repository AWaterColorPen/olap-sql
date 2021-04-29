// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api/proto/request.proto

package proto

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type METRIC_TYPE int32

const (
	METRIC_TYPE_METRIC_TYPE_UNKNOWN        METRIC_TYPE = 0
	METRIC_TYPE_METRIC_TYPE_COUNT          METRIC_TYPE = 1
	METRIC_TYPE_METRIC_TYPE_DISTINCT_COUNT METRIC_TYPE = 2
	METRIC_TYPE_METRIC_TYPE_SUM            METRIC_TYPE = 3
	METRIC_TYPE_METRIC_TYPE_POST           METRIC_TYPE = 4
	METRIC_TYPE_METRIC_TYPE_EXTENSION      METRIC_TYPE = 10
)

var METRIC_TYPE_name = map[int32]string{
	0:  "METRIC_TYPE_UNKNOWN",
	1:  "METRIC_TYPE_COUNT",
	2:  "METRIC_TYPE_DISTINCT_COUNT",
	3:  "METRIC_TYPE_SUM",
	4:  "METRIC_TYPE_POST",
	10: "METRIC_TYPE_EXTENSION",
}

var METRIC_TYPE_value = map[string]int32{
	"METRIC_TYPE_UNKNOWN":        0,
	"METRIC_TYPE_COUNT":          1,
	"METRIC_TYPE_DISTINCT_COUNT": 2,
	"METRIC_TYPE_SUM":            3,
	"METRIC_TYPE_POST":           4,
	"METRIC_TYPE_EXTENSION":      10,
}

func (x METRIC_TYPE) String() string {
	return proto.EnumName(METRIC_TYPE_name, int32(x))
}

func (METRIC_TYPE) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_336e959edb61afc4, []int{0}
}

type FILTER_OPERATOR_TYPE int32

const (
	FILTER_OPERATOR_TYPE_FILTER_OPERATOR_UNKNOWN        FILTER_OPERATOR_TYPE = 0
	FILTER_OPERATOR_TYPE_FILTER_OPERATOR_IN             FILTER_OPERATOR_TYPE = 1
	FILTER_OPERATOR_TYPE_FILTER_OPERATOR_NOT_IN         FILTER_OPERATOR_TYPE = 2
	FILTER_OPERATOR_TYPE_FILTER_OPERATOR_LESS_EQUALS    FILTER_OPERATOR_TYPE = 3
	FILTER_OPERATOR_TYPE_FILTER_OPERATOR_LESS           FILTER_OPERATOR_TYPE = 4
	FILTER_OPERATOR_TYPE_FILTER_OPERATOR_GREATER_EQUALS FILTER_OPERATOR_TYPE = 5
	FILTER_OPERATOR_TYPE_FILTER_OPERATOR_GREATER        FILTER_OPERATOR_TYPE = 6
	FILTER_OPERATOR_TYPE_FILTER_OPERATOR_LIKE           FILTER_OPERATOR_TYPE = 7
	FILTER_OPERATOR_TYPE_FILTER_OPERATOR_EXTENSION      FILTER_OPERATOR_TYPE = 20
)

var FILTER_OPERATOR_TYPE_name = map[int32]string{
	0:  "FILTER_OPERATOR_UNKNOWN",
	1:  "FILTER_OPERATOR_IN",
	2:  "FILTER_OPERATOR_NOT_IN",
	3:  "FILTER_OPERATOR_LESS_EQUALS",
	4:  "FILTER_OPERATOR_LESS",
	5:  "FILTER_OPERATOR_GREATER_EQUALS",
	6:  "FILTER_OPERATOR_GREATER",
	7:  "FILTER_OPERATOR_LIKE",
	20: "FILTER_OPERATOR_EXTENSION",
}

var FILTER_OPERATOR_TYPE_value = map[string]int32{
	"FILTER_OPERATOR_UNKNOWN":        0,
	"FILTER_OPERATOR_IN":             1,
	"FILTER_OPERATOR_NOT_IN":         2,
	"FILTER_OPERATOR_LESS_EQUALS":    3,
	"FILTER_OPERATOR_LESS":           4,
	"FILTER_OPERATOR_GREATER_EQUALS": 5,
	"FILTER_OPERATOR_GREATER":        6,
	"FILTER_OPERATOR_LIKE":           7,
	"FILTER_OPERATOR_EXTENSION":      20,
}

func (x FILTER_OPERATOR_TYPE) String() string {
	return proto.EnumName(FILTER_OPERATOR_TYPE_name, int32(x))
}

func (FILTER_OPERATOR_TYPE) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_336e959edb61afc4, []int{1}
}

type VALUE_TYPE int32

const (
	VALUE_TYPE_VALUE_UNKNOWN VALUE_TYPE = 0
	VALUE_TYPE_VALUE_STRING  VALUE_TYPE = 1
	VALUE_TYPE_VALUE_INTEGER VALUE_TYPE = 2
	VALUE_TYPE_VALUE_FLOAT   VALUE_TYPE = 3
)

var VALUE_TYPE_name = map[int32]string{
	0: "VALUE_UNKNOWN",
	1: "VALUE_STRING",
	2: "VALUE_INTEGER",
	3: "VALUE_FLOAT",
}

var VALUE_TYPE_value = map[string]int32{
	"VALUE_UNKNOWN": 0,
	"VALUE_STRING":  1,
	"VALUE_INTEGER": 2,
	"VALUE_FLOAT":   3,
}

func (x VALUE_TYPE) String() string {
	return proto.EnumName(VALUE_TYPE_name, int32(x))
}

func (VALUE_TYPE) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_336e959edb61afc4, []int{2}
}

type DATA_SOURCE_TYPE int32

const (
	DATA_SOURCE_TYPE_DATA_SOURCE_UNKNOWN    DATA_SOURCE_TYPE = 0
	DATA_SOURCE_TYPE_DATA_SOURCE_CLICKHOUSE DATA_SOURCE_TYPE = 1
	DATA_SOURCE_TYPE_DATA_SOURCE_DRUID      DATA_SOURCE_TYPE = 2
	DATA_SOURCE_TYPE_DATA_SOURCE_KYLIN      DATA_SOURCE_TYPE = 3
	DATA_SOURCE_TYPE_DATA_SOURCE_PRESTO     DATA_SOURCE_TYPE = 4
)

var DATA_SOURCE_TYPE_name = map[int32]string{
	0: "DATA_SOURCE_UNKNOWN",
	1: "DATA_SOURCE_CLICKHOUSE",
	2: "DATA_SOURCE_DRUID",
	3: "DATA_SOURCE_KYLIN",
	4: "DATA_SOURCE_PRESTO",
}

var DATA_SOURCE_TYPE_value = map[string]int32{
	"DATA_SOURCE_UNKNOWN":    0,
	"DATA_SOURCE_CLICKHOUSE": 1,
	"DATA_SOURCE_DRUID":      2,
	"DATA_SOURCE_KYLIN":      3,
	"DATA_SOURCE_PRESTO":     4,
}

func (x DATA_SOURCE_TYPE) String() string {
	return proto.EnumName(DATA_SOURCE_TYPE_name, int32(x))
}

func (DATA_SOURCE_TYPE) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_336e959edb61afc4, []int{3}
}

type Metric struct {
	Type                 METRIC_TYPE `protobuf:"varint,1,opt,name=type,proto3,enum=proto.METRIC_TYPE" json:"type,omitempty"`
	FieldName            string      `protobuf:"bytes,2,opt,name=field_name,json=fieldName,proto3" json:"field_name,omitempty"`
	Name                 string      `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	ExtensionValue       string      `protobuf:"bytes,4,opt,name=extension_value,json=extensionValue,proto3" json:"extension_value,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *Metric) Reset()         { *m = Metric{} }
func (m *Metric) String() string { return proto.CompactTextString(m) }
func (*Metric) ProtoMessage()    {}
func (*Metric) Descriptor() ([]byte, []int) {
	return fileDescriptor_336e959edb61afc4, []int{0}
}

func (m *Metric) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Metric.Unmarshal(m, b)
}
func (m *Metric) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Metric.Marshal(b, m, deterministic)
}
func (m *Metric) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Metric.Merge(m, src)
}
func (m *Metric) XXX_Size() int {
	return xxx_messageInfo_Metric.Size(m)
}
func (m *Metric) XXX_DiscardUnknown() {
	xxx_messageInfo_Metric.DiscardUnknown(m)
}

var xxx_messageInfo_Metric proto.InternalMessageInfo

func (m *Metric) GetType() METRIC_TYPE {
	if m != nil {
		return m.Type
	}
	return METRIC_TYPE_METRIC_TYPE_UNKNOWN
}

func (m *Metric) GetFieldName() string {
	if m != nil {
		return m.FieldName
	}
	return ""
}

func (m *Metric) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Metric) GetExtensionValue() string {
	if m != nil {
		return m.ExtensionValue
	}
	return ""
}

type Dimension struct {
	FieldName            string   `protobuf:"bytes,1,opt,name=field_name,json=fieldName,proto3" json:"field_name,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Dimension) Reset()         { *m = Dimension{} }
func (m *Dimension) String() string { return proto.CompactTextString(m) }
func (*Dimension) ProtoMessage()    {}
func (*Dimension) Descriptor() ([]byte, []int) {
	return fileDescriptor_336e959edb61afc4, []int{1}
}

func (m *Dimension) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Dimension.Unmarshal(m, b)
}
func (m *Dimension) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Dimension.Marshal(b, m, deterministic)
}
func (m *Dimension) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Dimension.Merge(m, src)
}
func (m *Dimension) XXX_Size() int {
	return xxx_messageInfo_Dimension.Size(m)
}
func (m *Dimension) XXX_DiscardUnknown() {
	xxx_messageInfo_Dimension.DiscardUnknown(m)
}

var xxx_messageInfo_Dimension proto.InternalMessageInfo

func (m *Dimension) GetFieldName() string {
	if m != nil {
		return m.FieldName
	}
	return ""
}

func (m *Dimension) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type Filter struct {
	OperatorType         FILTER_OPERATOR_TYPE `protobuf:"varint,1,opt,name=operator_type,json=operatorType,proto3,enum=proto.FILTER_OPERATOR_TYPE" json:"operator_type,omitempty"`
	ValueType            VALUE_TYPE           `protobuf:"varint,2,opt,name=value_type,json=valueType,proto3,enum=proto.VALUE_TYPE" json:"value_type,omitempty"`
	Name                 string               `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Value                []string             `protobuf:"bytes,4,rep,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *Filter) Reset()         { *m = Filter{} }
func (m *Filter) String() string { return proto.CompactTextString(m) }
func (*Filter) ProtoMessage()    {}
func (*Filter) Descriptor() ([]byte, []int) {
	return fileDescriptor_336e959edb61afc4, []int{2}
}

func (m *Filter) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Filter.Unmarshal(m, b)
}
func (m *Filter) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Filter.Marshal(b, m, deterministic)
}
func (m *Filter) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Filter.Merge(m, src)
}
func (m *Filter) XXX_Size() int {
	return xxx_messageInfo_Filter.Size(m)
}
func (m *Filter) XXX_DiscardUnknown() {
	xxx_messageInfo_Filter.DiscardUnknown(m)
}

var xxx_messageInfo_Filter proto.InternalMessageInfo

func (m *Filter) GetOperatorType() FILTER_OPERATOR_TYPE {
	if m != nil {
		return m.OperatorType
	}
	return FILTER_OPERATOR_TYPE_FILTER_OPERATOR_UNKNOWN
}

func (m *Filter) GetValueType() VALUE_TYPE {
	if m != nil {
		return m.ValueType
	}
	return VALUE_TYPE_VALUE_UNKNOWN
}

func (m *Filter) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Filter) GetValue() []string {
	if m != nil {
		return m.Value
	}
	return nil
}

type JoinOn struct {
	Key1                 string   `protobuf:"bytes,1,opt,name=key1,proto3" json:"key1,omitempty"`
	Key2                 string   `protobuf:"bytes,2,opt,name=key2,proto3" json:"key2,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *JoinOn) Reset()         { *m = JoinOn{} }
func (m *JoinOn) String() string { return proto.CompactTextString(m) }
func (*JoinOn) ProtoMessage()    {}
func (*JoinOn) Descriptor() ([]byte, []int) {
	return fileDescriptor_336e959edb61afc4, []int{3}
}

func (m *JoinOn) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_JoinOn.Unmarshal(m, b)
}
func (m *JoinOn) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_JoinOn.Marshal(b, m, deterministic)
}
func (m *JoinOn) XXX_Merge(src proto.Message) {
	xxx_messageInfo_JoinOn.Merge(m, src)
}
func (m *JoinOn) XXX_Size() int {
	return xxx_messageInfo_JoinOn.Size(m)
}
func (m *JoinOn) XXX_DiscardUnknown() {
	xxx_messageInfo_JoinOn.DiscardUnknown(m)
}

var xxx_messageInfo_JoinOn proto.InternalMessageInfo

func (m *JoinOn) GetKey1() string {
	if m != nil {
		return m.Key1
	}
	return ""
}

func (m *JoinOn) GetKey2() string {
	if m != nil {
		return m.Key2
	}
	return ""
}

type Join struct {
	Table1               string    `protobuf:"bytes,1,opt,name=table1,proto3" json:"table1,omitempty"`
	Table2               string    `protobuf:"bytes,2,opt,name=table2,proto3" json:"table2,omitempty"`
	On                   []*JoinOn `protobuf:"bytes,3,rep,name=on,proto3" json:"on,omitempty"`
	Filters              []*Filter `protobuf:"bytes,4,rep,name=filters,proto3" json:"filters,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *Join) Reset()         { *m = Join{} }
func (m *Join) String() string { return proto.CompactTextString(m) }
func (*Join) ProtoMessage()    {}
func (*Join) Descriptor() ([]byte, []int) {
	return fileDescriptor_336e959edb61afc4, []int{4}
}

func (m *Join) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Join.Unmarshal(m, b)
}
func (m *Join) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Join.Marshal(b, m, deterministic)
}
func (m *Join) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Join.Merge(m, src)
}
func (m *Join) XXX_Size() int {
	return xxx_messageInfo_Join.Size(m)
}
func (m *Join) XXX_DiscardUnknown() {
	xxx_messageInfo_Join.DiscardUnknown(m)
}

var xxx_messageInfo_Join proto.InternalMessageInfo

func (m *Join) GetTable1() string {
	if m != nil {
		return m.Table1
	}
	return ""
}

func (m *Join) GetTable2() string {
	if m != nil {
		return m.Table2
	}
	return ""
}

func (m *Join) GetOn() []*JoinOn {
	if m != nil {
		return m.On
	}
	return nil
}

func (m *Join) GetFilters() []*Filter {
	if m != nil {
		return m.Filters
	}
	return nil
}

type DataSource struct {
	Type                 DATA_SOURCE_TYPE `protobuf:"varint,1,opt,name=type,proto3,enum=proto.DATA_SOURCE_TYPE" json:"type,omitempty"`
	Name                 string           `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	SubRequest           *Request         `protobuf:"bytes,3,opt,name=sub_request,json=subRequest,proto3" json:"sub_request,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *DataSource) Reset()         { *m = DataSource{} }
func (m *DataSource) String() string { return proto.CompactTextString(m) }
func (*DataSource) ProtoMessage()    {}
func (*DataSource) Descriptor() ([]byte, []int) {
	return fileDescriptor_336e959edb61afc4, []int{5}
}

func (m *DataSource) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DataSource.Unmarshal(m, b)
}
func (m *DataSource) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DataSource.Marshal(b, m, deterministic)
}
func (m *DataSource) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DataSource.Merge(m, src)
}
func (m *DataSource) XXX_Size() int {
	return xxx_messageInfo_DataSource.Size(m)
}
func (m *DataSource) XXX_DiscardUnknown() {
	xxx_messageInfo_DataSource.DiscardUnknown(m)
}

var xxx_messageInfo_DataSource proto.InternalMessageInfo

func (m *DataSource) GetType() DATA_SOURCE_TYPE {
	if m != nil {
		return m.Type
	}
	return DATA_SOURCE_TYPE_DATA_SOURCE_UNKNOWN
}

func (m *DataSource) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *DataSource) GetSubRequest() *Request {
	if m != nil {
		return m.SubRequest
	}
	return nil
}

type Request struct {
	Metrics              []*Metric    `protobuf:"bytes,1,rep,name=metrics,proto3" json:"metrics,omitempty"`
	Dimensions           []*Dimension `protobuf:"bytes,2,rep,name=dimensions,proto3" json:"dimensions,omitempty"`
	Filters              []*Filter    `protobuf:"bytes,3,rep,name=filters,proto3" json:"filters,omitempty"`
	Joins                []*Join      `protobuf:"bytes,4,rep,name=joins,proto3" json:"joins,omitempty"`
	DataSource           *DataSource  `protobuf:"bytes,5,opt,name=data_source,json=dataSource,proto3" json:"data_source,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *Request) Reset()         { *m = Request{} }
func (m *Request) String() string { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()    {}
func (*Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_336e959edb61afc4, []int{6}
}

func (m *Request) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Request.Unmarshal(m, b)
}
func (m *Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Request.Marshal(b, m, deterministic)
}
func (m *Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Request.Merge(m, src)
}
func (m *Request) XXX_Size() int {
	return xxx_messageInfo_Request.Size(m)
}
func (m *Request) XXX_DiscardUnknown() {
	xxx_messageInfo_Request.DiscardUnknown(m)
}

var xxx_messageInfo_Request proto.InternalMessageInfo

func (m *Request) GetMetrics() []*Metric {
	if m != nil {
		return m.Metrics
	}
	return nil
}

func (m *Request) GetDimensions() []*Dimension {
	if m != nil {
		return m.Dimensions
	}
	return nil
}

func (m *Request) GetFilters() []*Filter {
	if m != nil {
		return m.Filters
	}
	return nil
}

func (m *Request) GetJoins() []*Join {
	if m != nil {
		return m.Joins
	}
	return nil
}

func (m *Request) GetDataSource() *DataSource {
	if m != nil {
		return m.DataSource
	}
	return nil
}

func init() {
	proto.RegisterEnum("proto.METRIC_TYPE", METRIC_TYPE_name, METRIC_TYPE_value)
	proto.RegisterEnum("proto.FILTER_OPERATOR_TYPE", FILTER_OPERATOR_TYPE_name, FILTER_OPERATOR_TYPE_value)
	proto.RegisterEnum("proto.VALUE_TYPE", VALUE_TYPE_name, VALUE_TYPE_value)
	proto.RegisterEnum("proto.DATA_SOURCE_TYPE", DATA_SOURCE_TYPE_name, DATA_SOURCE_TYPE_value)
	proto.RegisterType((*Metric)(nil), "proto.Metric")
	proto.RegisterType((*Dimension)(nil), "proto.Dimension")
	proto.RegisterType((*Filter)(nil), "proto.Filter")
	proto.RegisterType((*JoinOn)(nil), "proto.JoinOn")
	proto.RegisterType((*Join)(nil), "proto.Join")
	proto.RegisterType((*DataSource)(nil), "proto.DataSource")
	proto.RegisterType((*Request)(nil), "proto.Request")
}

func init() { proto.RegisterFile("api/proto/request.proto", fileDescriptor_336e959edb61afc4) }

var fileDescriptor_336e959edb61afc4 = []byte{
	// 795 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x55, 0xd1, 0x72, 0xdb, 0x44,
	0x14, 0x45, 0x92, 0xed, 0x8c, 0xaf, 0x9b, 0x64, 0xb3, 0x4d, 0x13, 0x35, 0x99, 0x94, 0xa0, 0x07,
	0xea, 0x31, 0xd3, 0xb8, 0x98, 0x77, 0x06, 0x63, 0x6f, 0x82, 0xb0, 0x23, 0x99, 0xd5, 0xaa, 0x50,
	0x5e, 0x34, 0xb2, 0xbd, 0x05, 0x81, 0x2d, 0xb9, 0x92, 0x0c, 0xf4, 0xa5, 0x5f, 0x00, 0x7f, 0xc0,
	0x1b, 0x33, 0x7c, 0x19, 0x1f, 0xc2, 0x68, 0x25, 0x59, 0x5b, 0x61, 0xfa, 0xa4, 0x7b, 0xcf, 0xb9,
	0x47, 0x7b, 0xef, 0xd1, 0x5d, 0x1b, 0xce, 0xfd, 0x4d, 0xd0, 0xdf, 0xc4, 0x51, 0x1a, 0xf5, 0x63,
	0xfe, 0x7a, 0xcb, 0x93, 0xf4, 0x46, 0x64, 0xb8, 0x29, 0x1e, 0xc6, 0xef, 0x0a, 0xb4, 0xee, 0x79,
	0x1a, 0x07, 0x0b, 0xfc, 0x31, 0x34, 0xd2, 0x37, 0x1b, 0xae, 0x2b, 0xd7, 0x4a, 0xf7, 0x68, 0x80,
	0xf3, 0xba, 0x9b, 0x7b, 0xc2, 0xa8, 0x39, 0xf2, 0xd8, 0xcb, 0x19, 0xa1, 0x82, 0xc7, 0x57, 0x00,
	0xaf, 0x02, 0xbe, 0x5a, 0x7a, 0xa1, 0xbf, 0xe6, 0xba, 0x7a, 0xad, 0x74, 0xdb, 0xb4, 0x2d, 0x10,
	0xcb, 0x5f, 0x73, 0x8c, 0xa1, 0x21, 0x08, 0x4d, 0x10, 0x22, 0xc6, 0x4f, 0xe1, 0x98, 0xff, 0x96,
	0xf2, 0x30, 0x09, 0xa2, 0xd0, 0xfb, 0xc5, 0x5f, 0x6d, 0xb9, 0xde, 0x10, 0xf4, 0xd1, 0x0e, 0x7e,
	0x91, 0xa1, 0xc6, 0xe7, 0xd0, 0x1e, 0x07, 0xeb, 0x1c, 0xa9, 0x1d, 0xa4, 0xfc, 0xdf, 0x41, 0x6a,
	0x75, 0x90, 0xf1, 0xb7, 0x02, 0xad, 0xdb, 0x60, 0x95, 0xf2, 0x18, 0x7f, 0x01, 0x87, 0xd1, 0x86,
	0xc7, 0x7e, 0x1a, 0xc5, 0x9e, 0x34, 0xd7, 0x65, 0x31, 0xd7, 0xad, 0x39, 0x65, 0x84, 0x7a, 0xf6,
	0x8c, 0xd0, 0x21, 0xb3, 0x69, 0x3e, 0xe0, 0x83, 0x52, 0xc1, 0xb2, 0x41, 0x9f, 0x03, 0x88, 0x5e,
	0x73, 0xb9, 0x2a, 0xe4, 0x27, 0x85, 0xfc, 0xc5, 0x70, 0xea, 0x92, 0x5c, 0xd4, 0x16, 0x45, 0x42,
	0xb1, 0x6f, 0xf6, 0x53, 0x68, 0x96, 0x13, 0x6b, 0xdd, 0x36, 0xcd, 0x13, 0xe3, 0x39, 0xb4, 0xbe,
	0x8e, 0x82, 0xd0, 0x0e, 0x33, 0xcd, 0xcf, 0xfc, 0xcd, 0xa7, 0xc5, 0x7c, 0x22, 0x2e, 0xb0, 0x41,
	0x39, 0x5a, 0x16, 0x1b, 0x6f, 0xa1, 0x91, 0x29, 0xf0, 0x19, 0xb4, 0x52, 0x7f, 0xbe, 0xe2, 0xa5,
	0xa2, 0xc8, 0x76, 0x78, 0xa9, 0x2a, 0x32, 0x7c, 0x05, 0x6a, 0x14, 0xea, 0xda, 0xb5, 0xd6, 0xed,
	0x0c, 0x0e, 0x8b, 0xee, 0xf3, 0xa3, 0xa9, 0x1a, 0x85, 0xf8, 0x29, 0x1c, 0xbc, 0x12, 0x86, 0x25,
	0xa2, 0xc1, 0xaa, 0x26, 0xb7, 0x91, 0x96, 0xac, 0xf1, 0x16, 0x60, 0xec, 0xa7, 0xbe, 0x13, 0x6d,
	0xe3, 0x05, 0xc7, 0x9f, 0xbc, 0xb3, 0x2c, 0xe7, 0x85, 0x66, 0x3c, 0x64, 0x43, 0xcf, 0xb1, 0x5d,
	0x3a, 0x22, 0xf2, 0xc6, 0xec, 0xf9, 0x52, 0xb8, 0x0f, 0x9d, 0x64, 0x3b, 0xf7, 0x8a, 0xa5, 0x14,
	0x8e, 0x75, 0x06, 0x47, 0xc5, 0x7b, 0x68, 0x8e, 0x52, 0x48, 0xb6, 0xf3, 0x22, 0x36, 0xfe, 0x51,
	0xe0, 0xa0, 0x88, 0xb3, 0xa6, 0xd7, 0x62, 0x69, 0x13, 0x5d, 0x79, 0xa7, 0xe9, 0x7c, 0x95, 0x69,
	0xc9, 0x66, 0x9f, 0x70, 0x59, 0xee, 0x53, 0xa2, 0xab, 0xa2, 0x16, 0x95, 0xcd, 0x96, 0x04, 0x95,
	0x6a, 0x64, 0x3f, 0xb4, 0xf7, 0xf9, 0x81, 0x3f, 0x82, 0xe6, 0x4f, 0x51, 0x10, 0x96, 0xb6, 0x75,
	0x24, 0x6b, 0x69, 0xce, 0xe0, 0x01, 0x74, 0x96, 0x7e, 0xea, 0x7b, 0x89, 0xf0, 0x4c, 0x6f, 0x8a,
	0x19, 0xcb, 0x0d, 0xaa, 0xcc, 0xa4, 0xb0, 0xdc, 0xc5, 0xbd, 0xbf, 0x14, 0xe8, 0x48, 0x77, 0x0e,
	0x9f, 0xc3, 0x43, 0x29, 0xf5, 0x5c, 0x6b, 0x62, 0xd9, 0xdf, 0x5a, 0xe8, 0x03, 0xfc, 0x08, 0x4e,
	0x64, 0x62, 0x64, 0xbb, 0x16, 0x43, 0x0a, 0x7e, 0x02, 0x17, 0x32, 0x3c, 0x36, 0x1d, 0x66, 0x5a,
	0x23, 0x56, 0xf0, 0x2a, 0x7e, 0x08, 0xc7, 0x32, 0xef, 0xb8, 0xf7, 0x48, 0xc3, 0xa7, 0x80, 0x64,
	0x70, 0x66, 0x3b, 0x0c, 0x35, 0xf0, 0x63, 0x78, 0x24, 0xa3, 0xe4, 0x3b, 0x46, 0x2c, 0xc7, 0xb4,
	0x2d, 0x04, 0xbd, 0x3f, 0x55, 0x38, 0xdd, 0x77, 0x83, 0xf0, 0x25, 0x9c, 0xd7, 0xf1, 0xaa, 0xe5,
	0x33, 0xc0, 0x75, 0xd2, 0xb4, 0x90, 0x82, 0x2f, 0xe0, 0xac, 0x8e, 0x5b, 0x36, 0xcb, 0x38, 0x15,
	0x7f, 0x08, 0x97, 0x75, 0x6e, 0x4a, 0x1c, 0xc7, 0x23, 0xdf, 0xb8, 0xc3, 0xa9, 0x83, 0x34, 0xac,
	0xff, 0xb7, 0x93, 0xac, 0x00, 0x35, 0xb0, 0x01, 0x4f, 0xea, 0xcc, 0x1d, 0x25, 0xc3, 0x0c, 0x28,
	0xd4, 0xcd, 0x7d, 0xfd, 0x16, 0x35, 0xa8, 0xb5, 0xf7, 0xd5, 0xe6, 0x84, 0xa0, 0x03, 0x7c, 0x05,
	0x8f, 0xeb, 0x4c, 0x65, 0xcf, 0x69, 0xcf, 0x05, 0xa8, 0x7e, 0x20, 0xf0, 0x09, 0x1c, 0xe6, 0x59,
	0xe5, 0x04, 0x82, 0x07, 0x39, 0xe4, 0x30, 0x6a, 0x5a, 0x77, 0x48, 0xa9, 0x8a, 0x4c, 0x8b, 0x91,
	0x3b, 0x42, 0x91, 0x8a, 0x8f, 0xa1, 0x93, 0x43, 0xb7, 0x53, 0x7b, 0xc8, 0x90, 0xd6, 0xfb, 0x43,
	0x01, 0x54, 0xbf, 0x62, 0xd9, 0x82, 0xc8, 0x58, 0x75, 0xc6, 0x05, 0x9c, 0xc9, 0xc4, 0x68, 0x6a,
	0x8e, 0x26, 0x5f, 0xd9, 0xae, 0x43, 0x90, 0x92, 0x2d, 0x8f, 0xcc, 0x8d, 0xa9, 0x6b, 0x8e, 0x91,
	0x5a, 0x87, 0x27, 0x2f, 0xa7, 0xa6, 0x85, 0xb4, 0xec, 0xbb, 0xc9, 0xf0, 0x8c, 0x12, 0x87, 0xd9,
	0xa8, 0xf1, 0x65, 0xff, 0xfb, 0x67, 0xd7, 0x3f, 0x04, 0xe9, 0x8f, 0xdb, 0xf9, 0xcd, 0x22, 0x5a,
	0xf7, 0xfd, 0x5f, 0xfd, 0x94, 0xc7, 0x8b, 0x68, 0x15, 0xc5, 0x1b, 0x1e, 0xf6, 0xa3, 0x95, 0xbf,
	0x79, 0x96, 0xbc, 0x5e, 0xf5, 0x77, 0x7f, 0x41, 0xf3, 0x96, 0x78, 0x7c, 0xf6, 0x6f, 0x00, 0x00,
	0x00, 0xff, 0xff, 0x24, 0x7a, 0xb2, 0xa6, 0x96, 0x06, 0x00, 0x00,
}