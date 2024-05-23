package bean

import (
	"github.com/gin-gonic/gin"
	"github.com/simonalong/gole/logger"
	"github.com/simonalong/gole/server/rsp"
	"github.com/simonalong/gole/util"
	"reflect"
	"strings"
)

var BeanMap map[string]any

func init() {
	BeanMap = map[string]any{}
}

// AddBean 添加bean
// 强烈建议：bean使用指针
func AddBean(beanName string, beanPtr any) {
	BeanMap[beanName] = beanPtr
}

func GetBean(beanName string) any {
	if beanV, exit := BeanMap[beanName]; exit {
		return beanV
	}
	return nil
}

func Clean() {
	BeanMap = map[string]any{}
}

// 模糊搜索
func GetBeanNames(beanName string) []string {
	if beanName == "" {
		j := 0
		keys := make([]string, len(BeanMap))
		for k := range BeanMap {
			keys[j] = k
			j++
		}
		return keys
	} else {
		var keys []string
		for k := range BeanMap {
			if strings.Contains(k, beanName) {
				keys = append(keys, k)
			}
		}
		return keys
	}
}

// 前缀搜索
func GetBeanWithNamePre(beanName string) []any {
	if beanName == "" {
		return nil
	} else {
		var values []any
		for k, v := range BeanMap {
			if strings.HasPrefix(k, beanName) {
				values = append(values, v)
			}
		}
		return values
	}
}

func ExistBean(beanName string) bool {
	_, exist := BeanMap[beanName]
	return exist
}

// @parameterValueMap p1、p2、p3...这种表示的是第一个、第二个、第三个参数的值
func CallFun(beanName, methodName string, parameterValueMap map[string]any) []any {
	if beanValue, exist := BeanMap[beanName]; exist {
		fType := reflect.TypeOf(beanValue)
		var result []any
		for index, num := 0, fType.NumMethod(); index < num; index++ {
			method := fType.Method(index)

			// 私有字段不处理
			if !isc.IsPublic(method.Name) {
				continue
			}

			if method.Name != methodName {
				continue
			}

			parameterNum := method.Type.NumIn()
			var in []reflect.Value
			in = append(in, reflect.ValueOf(beanValue))
			for i := 1; i < parameterNum; i++ {
				if v, exit := parameterValueMap["p"+isc.ToString(i)]; exit {
					pV, err := isc.ToValue(v, method.Type.In(i).Kind())
					if err != nil {
						continue
					}
					in = append(in, reflect.ValueOf(pV))
				}
			}

			if len(in) != parameterNum {
				return nil
			}
			vs := method.Func.Call(in)
			for _, v := range vs {
				result = append(result, v.Interface())
			}
		}
		return result
	}
	return nil
}

func GetField(beanName, fieldName string) any {
	if beanValue, exist := BeanMap[beanName]; exist {
		fValue := reflect.ValueOf(beanValue)
		fType := reflect.TypeOf(beanValue)

		if fType.Kind() == reflect.Ptr {
			fType = fType.Elem()
			fValue = fValue.Elem()
		}

		for index, num := 0, fType.NumField(); index < num; index++ {
			field := fType.Field(index)

			// 私有字段不处理
			if !isc.IsPublic(field.Name) {
				continue
			}

			if field.Name == fieldName {
				return fValue.Field(index).Interface()
			}
		}
	}
	return nil
}

// SetField 修改属性的话，请将对象设置为指针，否则不生效
func SetField(beanName string, fieldName string, fieldValue any) {
	if beanValue, exist := BeanMap[beanName]; exist {
		// 私有字段不处理
		if !isc.IsPublic(fieldName) {
			return
		}

		fValue := reflect.ValueOf(beanValue)
		fType := reflect.TypeOf(beanValue)

		if fType.Kind() == reflect.Ptr {
			fType = fType.Elem()
			fValue = fValue.Elem()
		} else {
			logger.Warn("对象名【%v】对应的对象不是指针类型，无法修改属性的值", beanName)
			return
		}

		if _, exist := fType.FieldByName(fieldName); exist {
			v, err := isc.ToValue(fieldValue, fValue.FieldByName(fieldName).Kind())
			if err != nil {
				return
			}
			fValue.FieldByName(fieldName).Set(reflect.ValueOf(v))
		}
	}
}

func DebugBeanAll(c *gin.Context) {
	rsp.SuccessOfStandard(c, GetBeanNames(""))
}

func DebugBeanList(c *gin.Context) {
	rsp.SuccessOfStandard(c, GetBeanNames(c.Param("name")))
}

func DebugBeanGetField(c *gin.Context) {
	fieldGetReq := FieldGetReq{}
	err := isc.DataToObject(c.Request.Body, &fieldGetReq)
	if err != nil {
		return
	}
	rsp.SuccessOfStandard(c, GetField(fieldGetReq.Bean, fieldGetReq.Field))
}

func DebugBeanSetField(c *gin.Context) {
	fieldSetReq := FieldSetReq{}
	err := isc.DataToObject(c.Request.Body, &fieldSetReq)
	if err != nil {
		return
	}
	SetField(fieldSetReq.Bean, fieldSetReq.Field, fieldSetReq.Value)
	rsp.SuccessOfStandard(c, fieldSetReq.Value)
}

func DebugBeanFunCall(c *gin.Context) {
	funCallReq := FunCallReq{}
	err := isc.DataToObject(c.Request.Body, &funCallReq)
	if err != nil {
		return
	}
	rsp.SuccessOfStandard(c, CallFun(funCallReq.Bean, funCallReq.Fun, funCallReq.Parameter))
}

func BeanTest() {
	logger.Warn("test, ttt")
}

type FieldGetReq struct {
	Bean  string
	Field string
}

type FieldSetReq struct {
	Bean  string
	Field string
	Value any
}

type FunCallReq struct {
	Bean      string
	Fun       string
	Parameter map[string]any
}
