package util

import (
	"encoding/json"
	"errors"
	"fmt"
	cmap "github.com/orcaman/concurrent-map"
	goleTime "github.com/simonalong/gole/time"
	"reflect"
	"strings"
	"time"
)

/**
 * 提供新的map
 * 1. 提供类型转换
 * 2. 并发安全
 * 3. 提供有序性
 * 4. 提供与实体的转化功能
 */

type GoleMap struct {
	innerMap cmap.ConcurrentMap
	sort     bool
	keys     []string
}

func NewGoleMap() *GoleMap {
	return &GoleMap{
		innerMap: cmap.New(),
		sort:     false,
		keys:     make([]string, 0),
	}
}

func NewSortGoleMap() *GoleMap {
	return &GoleMap{
		innerMap: cmap.New(),
		sort:     true,
		keys:     make([]string, 0),
	}
}

// GoleMapOf 支持：k-v-k-v结构
// 默认无序，如果想要有序，请使用OfSort()
func GoleMapOf(parameters ...any) *GoleMap {
	if parameters == nil || len(parameters) == 0 {
		return NewGoleMap()
	}

	pGoleMap := NewGoleMap()
	for index := 0; index < len(parameters); index++ {
		data := parameters[index]
		if reflect.TypeOf(data).Kind() == reflect.String {
			key := data.(string)
			var value interface{}
			if (index + 1) < len(parameters) {
				value = parameters[index+1]
			}

			pGoleMap.Put(key, value)
			index++
		}
	}
	return pGoleMap
}

func GoleMapOfSort(parameters ...any) *GoleMap {
	if parameters == nil || len(parameters) == 0 {
		return NewSortGoleMap()
	}

	pGoleMap := NewSortGoleMap()
	for index := 0; index < len(parameters); index++ {
		data := parameters[index]
		if reflect.TypeOf(data).Kind() == reflect.String {
			key := data.(string)
			var value interface{}
			if (index + 1) < len(parameters) {
				value = parameters[index+1]
			}

			pGoleMap.Put(key, value)
			index++
		}
	}
	return pGoleMap
}

func FromToGoleMap(entity interface{}) (*GoleMap, error) {
	if entity == nil {
		return nil, nil
	}
	entityType := reflect.TypeOf(entity)
	if entityType.Kind() == reflect.Map {
		return FromMapToGoleMap(entity.(map[string]interface{})), nil
	} else if entityType.Kind() == reflect.Struct {
		return FromEntityToGoleMap(entity), nil
	} else if entityType.Kind() == reflect.String {
		return FromJsonToGoleMap(entity.(string))
	} else {
		return nil, errors.New(fmt.Sprintf("暂时不支持除了map、struct和string之外的其他类型：%v", entityType.Kind().String()))
	}
}

// FromEntityToGoleMap 从实体转换为map，默认转换为有序map
func FromEntityToGoleMap(entity interface{}) *GoleMap {
	if entity == nil {
		return NewGoleMap()
	}

	valType := reflect.TypeOf(entity)
	if valType.Kind() == reflect.Ptr {
		valType = valType.Elem()
	}
	if valType == reflect.TypeOf(GoleMap{}) {
		return entity.(*GoleMap)
	}

	objType := reflect.TypeOf(entity)
	// 只接收结构体类型
	if objType.Kind() != reflect.Struct {
		return nil
	}

	entityMap := NewSortGoleMap()

	objValue := reflect.ValueOf(entity)
	for fieldIndex, num := 0, objType.NumField(); fieldIndex < num; fieldIndex++ {
		field := objType.Field(fieldIndex)
		if !IsPublic(field.Name) {
			continue
		}

		columnName := getFinalColumnName(field)

		fieldValue := objValue.Field(fieldIndex)
		entityMap.Put(columnName, fieldValue.Interface())
	}
	return entityMap
}

// FromMapToGoleMap 从map转换为neoMap，默认转换为有序map
func FromMapToGoleMap(dataMap map[string]interface{}) *GoleMap {
	if dataMap == nil || len(dataMap) == 0 {
		return NewGoleMap()
	}
	resultMap := NewSortGoleMap()
	for key, val := range dataMap {
		resultMap.Put(key, val)
	}
	return resultMap
}

func FromJsonToGoleMap(jsonOfContent string) (*GoleMap, error) {
	if jsonOfContent == "" {
		return NewGoleMap(), nil
	}
	resultMap := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonOfContent), &resultMap)
	if err != nil {
		return nil, err
	}

	return FromMapToGoleMap(resultMap), nil
}

func (receiver *GoleMap) ToEntity(pEntity interface{}) error {
	if receiver == nil {
		return nil
	}
	if pEntity == nil {
		return errors.New("对象指针为nil")
	}

	return MapToObject(receiver.innerMap.Items(), pEntity)
}

func (receiver *GoleMap) ToMap() map[string]interface{} {
	if receiver == nil {
		return nil
	}
	return receiver.innerMap.Items()
}

func (receiver *GoleMap) ToJson() string {
	if receiver == nil {
		return "{}"
	}
	return ToJsonString(receiver.innerMap)
}

// ToJsonOfSort 输出key有序的json
func (receiver *GoleMap) ToJsonOfSort() string {
	if receiver == nil || receiver.IsEmpty() {
		return "{}"
	}
	var jsonResult string
	jsonResult += "{"

	var kvs []string
	for _, key := range receiver.Keys() {
		val, _ := receiver.Get(key)
		if val == nil {
			continue
		}
		valType := reflect.TypeOf(val)
		valValue := reflect.ValueOf(val)
		if valType.Kind() == reflect.Ptr {
			valType = valType.Elem()
			valValue = valValue.Elem()
		}
		if IsStringType(valType) {
			kvs = append(kvs, fmt.Sprintf("\"%v\":\"%v\"", key, valValue.Interface()))
		} else if IsNumberType(valType) || IsBoolType(valType) {
			kvs = append(kvs, fmt.Sprintf("\"%v\":%v", key, valValue.Interface()))
		} else if IsTimeType(valType) {
			kvs = append(kvs, fmt.Sprintf("\"%v\":\"%v\"", key, valValue.Interface()))
		} else if IsArrayType(valType) {
			kvs = append(kvs, fmt.Sprintf("\"%v\":%v", key, ToJsonString(valValue.Interface())))
		} else {
			var valJson string
			if valType == reflect.TypeOf(GoleMap{}) {
				value := val.(*GoleMap)
				valJson = value.ToJsonOfSort()
			} else if valType.Kind() == reflect.Map {
				valJson = ToJsonString(valValue.Interface())
			} else {
				valJson = FromEntityToGoleMap(valValue.Interface()).ToJsonOfSort()
			}
			if valJson == "{}" {
				continue
			}
			kvs = append(kvs, fmt.Sprintf("\"%v\":%v", key, valJson))
		}
	}

	jsonResult += strings.Join(kvs, ",")
	jsonResult += "}"
	return jsonResult
}

func (receiver *GoleMap) ToString() string {
	if receiver == nil {
		return ""
	}
	var keyValue []string
	for _, key := range receiver.Keys() {
		val, _ := receiver.GetString(key)
		keyValue = append(keyValue, "\""+key+"\":\""+val+"\"")
	}
	return "[" + strings.Join(keyValue, ",") + "]"
}

func (receiver *GoleMap) Keys() []string {
	if receiver == nil {
		return nil
	}
	if receiver.sort {
		return receiver.keys
	} else {
		return receiver.innerMap.Keys()
	}
}

// SetSort 设置map为有序或者无序map
// 注意：
//  1. 如果从无序变为有序，且之前已经有一些数据，则之前的数据顺序至此固定，后续的顺序就按照添加的顺序固定
//  2. 如果从有序变为无序，且之前已经有一些数据，则顺序就完全乱掉了
func (receiver *GoleMap) SetSort(sort bool) *GoleMap {
	if receiver == nil {
		return receiver
	}
	if !receiver.sort && sort {
		receiver.keys = receiver.innerMap.Keys()
	} else if receiver.sort && !sort {
		receiver.keys = make([]string, 0)
	}
	receiver.sort = sort
	return receiver
}

func (receiver *GoleMap) IsEmpty() bool {
	if receiver == nil {
		return true
	}
	return len(receiver.innerMap.Keys()) == 0
}

// GoleMapAllIsEmpty 所有数据都为空，则返回true
func GoleMapAllIsEmpty(dataMaps []*GoleMap) bool {
	if dataMaps == nil {
		return true
	}

	for _, dataMap := range dataMaps {
		if dataMap.IsUnEmpty() {
			return false
		}
	}
	return true
}

func (receiver *GoleMap) IsUnEmpty() bool {
	if receiver == nil {
		return false
	}
	return len(receiver.innerMap.Keys()) != 0
}

func (receiver *GoleMap) Clone() *GoleMap {
	if receiver == nil {
		return nil
	}
	cloneMap := &GoleMap{
		innerMap: cmap.New(),
		sort:     receiver.sort,
		keys:     make([]string, 0),
	}
	for _, key := range receiver.Keys() {
		val, _ := receiver.Get(key)
		cloneMap.Put(key, val)
	}
	return cloneMap
}

// AsDeepMap 将对象转换为可以使用test.name.single.age这样访问的map，不过这个neoMap不建议使用，建议只作为读取用
func (receiver *GoleMap) AsDeepMap() *GoleMap {
	if receiver == nil {
		return nil
	}
	mapFromProperties, _ := MapToProperties(receiver.ToMap())
	propertiesFromMap, _ := PropertiesToMap(mapFromProperties)

	deepMap := &GoleMap{
		innerMap: cmap.New(),
		sort:     false,
		keys:     make([]string, 0),
	}
	for key, val := range propertiesFromMap {
		deepMap.Put(key, val)
	}
	return deepMap
}

func (receiver *GoleMap) Put(key string, value interface{}) *GoleMap {
	if receiver == nil {
		return nil
	}
	if key == "" {
		return receiver
	}
	receiver.innerMap.Set(key, value)
	if receiver.sort {
		receiver.keys = append(receiver.keys, key)
	}
	return receiver
}

func (receiver *GoleMap) Contain(key string) bool {
	if receiver == nil {
		return false
	}
	if key == "" {
		return false
	}
	_, exit := receiver.innerMap.Get(key)
	return exit
}

func (receiver *GoleMap) Get(key string) (interface{}, bool) {
	if receiver == nil {
		return nil, false
	}
	if key == "" {
		return nil, false
	}
	return receiver.innerMap.Get(key)
}

func (receiver *GoleMap) GetGoleMap(key string) (*GoleMap, bool) {
	if receiver == nil {
		return nil, false
	}
	if key == "" {
		return nil, false
	}
	val, exist := receiver.innerMap.Get(key)
	if !exist {
		return nil, false
	}
	return val.(*GoleMap), true
}

func (receiver *GoleMap) GetInt(key string) (int, bool) {
	if receiver == nil {
		return 0, false
	}
	if key == "" {
		return 0, false
	}
	d, exit := receiver.innerMap.Get(key)
	if exit {
		return ToInt(d), true
	} else {
		return 0, false
	}
}

func (receiver *GoleMap) GetInt8(key string) (int8, bool) {
	if receiver == nil {
		return 0, false
	}
	if key == "" {
		return 0, false
	}
	d, exit := receiver.innerMap.Get(key)
	if exit {
		return ToInt8(d), true
	} else {
		return 0, false
	}
}

func (receiver *GoleMap) GetInt16(key string) (int16, bool) {
	if receiver == nil {
		return 0, false
	}
	if key == "" {
		return 0, false
	}
	d, exit := receiver.innerMap.Get(key)
	if exit {
		return ToInt16(d), true
	} else {
		return 0, false
	}
}

func (receiver *GoleMap) GetInt32(key string) (int32, bool) {
	if receiver == nil {
		return 0, false
	}
	if key == "" {
		return 0, false
	}
	d, exit := receiver.innerMap.Get(key)
	if exit {
		return ToInt32(d), true
	} else {
		return 0, false
	}
}

func (receiver *GoleMap) GetInt64(key string) (int64, bool) {
	if receiver == nil {
		return 0, false
	}
	if key == "" {
		return 0, false
	}
	d, exit := receiver.innerMap.Get(key)
	if exit {
		return ToInt64(d), true
	} else {
		return 0, false
	}
}

func (receiver *GoleMap) GetUInt(key string) (uint, bool) {
	if receiver == nil {
		return 0, false
	}
	if key == "" {
		return 0, false
	}
	d, exit := receiver.innerMap.Get(key)
	if exit {
		return ToUInt(d), true
	} else {
		return 0, false
	}
}

func (receiver *GoleMap) GetUInt8(key string) (uint8, bool) {
	if receiver == nil {
		return 0, false
	}
	if key == "" {
		return 0, false
	}
	d, exit := receiver.innerMap.Get(key)
	if exit {
		return ToUInt8(d), true
	} else {
		return 0, false
	}
}

func (receiver *GoleMap) GetUInt16(key string) (uint16, bool) {
	if receiver == nil {
		return 0, false
	}
	if key == "" {
		return 0, false
	}
	d, exit := receiver.innerMap.Get(key)
	if exit {
		return ToUInt16(d), true
	} else {
		return 0, false
	}
}

func (receiver *GoleMap) GetUInt32(key string) (uint32, bool) {
	if receiver == nil {
		return 0, false
	}
	if key == "" {
		return 0, false
	}
	d, exit := receiver.innerMap.Get(key)
	if exit {
		return ToUInt32(d), true
	} else {
		return 0, false
	}
}

func (receiver *GoleMap) GetUInt64(key string) (uint64, bool) {
	if receiver == nil {
		return 0, false
	}
	if key == "" {
		return 0, false
	}
	d, exit := receiver.innerMap.Get(key)
	if exit {
		return ToUInt64(d), true
	} else {
		return 0, false
	}
}

func (receiver *GoleMap) GetFloat32(key string) (float32, bool) {
	if receiver == nil {
		return 0, false
	}
	if key == "" {
		return 0, false
	}
	d, exit := receiver.innerMap.Get(key)
	if exit {
		return ToFloat32(d), true
	} else {
		return 0, false
	}
}

func (receiver *GoleMap) GetFloat64(key string) (float64, bool) {
	if receiver == nil {
		return 0, false
	}
	if key == "" {
		return 0, false
	}
	d, exit := receiver.innerMap.Get(key)
	if exit {
		return ToFloat64(d), true
	} else {
		return 0, false
	}
}

func (receiver *GoleMap) GetBool(key string) (bool, bool) {
	if receiver == nil {
		return false, false
	}
	if key == "" {
		return false, false
	}
	d, exit := receiver.innerMap.Get(key)
	if exit {
		return ToBool(d), true
	} else {
		return false, false
	}
}

func (receiver *GoleMap) GetComplex64(key string) (complex64, bool) {
	if receiver == nil {
		return 0, false
	}
	if key == "" {
		return 0, false
	}
	d, exit := receiver.innerMap.Get(key)
	if exit {
		return ToComplex64(d), true
	} else {
		return 0, false
	}
}

func (receiver *GoleMap) GetComplex128(key string) (complex128, bool) {
	if receiver == nil {
		return 0, false
	}
	if key == "" {
		return 0, false
	}
	d, exit := receiver.innerMap.Get(key)
	if exit {
		return ToComplex128(d), true
	} else {
		return 0, false
	}
}

func (receiver *GoleMap) GetString(key string) (string, bool) {
	if receiver == nil {
		return "", false
	}
	if key == "" {
		return "", false
	}
	val, exit := receiver.innerMap.Get(key)
	if exit {
		if timestampVal, ok := val.(time.Time); ok {
			return goleTime.TimeToStringYmdHmsS(timestampVal), true
		} else {
			return ToString(val), true
		}
	} else {
		return "", false
	}
}

func (receiver *GoleMap) GetTime(key string) (time.Time, bool) {
	if receiver == nil {
		return time.Time{}, false
	}
	if key == "" {
		return time.Time{}, false
	}
	d, exit := receiver.innerMap.Get(key)
	if exit {
		return d.(time.Time), true
	} else {
		return time.Now(), false
	}
}

func (receiver *GoleMap) GetBytes(key string) ([]byte, bool) {
	if receiver == nil {
		return nil, false
	}
	if key == "" {
		return make([]byte, 0), false
	}
	d, exit := receiver.innerMap.Get(key)
	if exit {
		return []byte(ToString(d)), true
	} else {
		return nil, false
	}
}

func (receiver *GoleMap) Remove(key string) {
	if receiver == nil {
		return
	}
	receiver.innerMap.Remove(key)
	if receiver.sort {
		id := IndexOf(receiver.keys, key)
		receiver.keys = append(receiver.keys[:id], receiver.keys[id+1:]...)
	}
}
func (receiver *GoleMap) RemoveAll() {
	if receiver == nil {
		return
	}
	receiver.innerMap.Clear()
	if receiver.sort {
		receiver.keys = make([]string, 0)
	}
}
func (receiver *GoleMap) Clear() {
	if receiver == nil {
		return
	}
	receiver.innerMap.Clear()
	if receiver.sort {
		receiver.keys = make([]string, 0)
	}
}
func (receiver *GoleMap) Size() int {
	if receiver == nil {
		return 0
	}
	return len(receiver.innerMap.Keys())
}

func getFinalColumnName(field reflect.StructField) string {
	// 先读取标签column
	columnName := field.Tag.Get("column")
	if len(columnName) != 0 {
		return columnName
	}

	// 如果没有配置column标签，也可以使用json标签，这里也支持
	aliasJson := field.Tag.Get("json")
	if len(aliasJson) != 0 {
		index := strings.Index(aliasJson, ",")
		if index != -1 {
			return aliasJson[:index]
		} else {
			return aliasJson
		}
	}

	// 如果也没有配置json标签，则使用属性的属性名，将首字母变小写
	return ToLowerFirstPrefix(field.Name)
}
