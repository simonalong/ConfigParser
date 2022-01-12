package config

import (
	"flag"
	"fmt"
	"github.com/simonalong/gole/util"
	"github.com/simonalong/gole/yaml"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"
)

var appProperty *ApplicationProperty

// LoadConfig 默认读取./resources/下面的配置文件
// 支持yml、yaml、json、properties格式
// 优先级yaml > yml > properties > json
func LoadConfig() {
	LoadConfigWithRelativePath("./resources/")
}

// LoadConfigWithRelativePath 加载相对文件路径，相对路径是相对系统启动的位置部分
func LoadConfigWithRelativePath(resourceAbsPath string) {
	dir, _ := os.Getwd()
	pkg := strings.Replace(dir, "\\", "/", -1)
	LoadConfigWithAbsPath(path.Join(pkg, "", resourceAbsPath))
}

// LoadConfigWithAbsPath 加载资源文件目录的绝对路径内容，比如：/user/xxx/mmm-biz-service/resources/
// 支持yml、yaml、json、properties格式
// 优先级yaml > yml > properties > json
// 支持命令行：--app.profile xxx
func LoadConfigWithAbsPath(resourceAbsPath string) {
	if !strings.HasSuffix(resourceAbsPath, "/") {
		resourceAbsPath += "/"
	}
	files, err := ioutil.ReadDir(resourceAbsPath)
	if err != nil {
		fmt.Printf("read fail, resource: %v, err %v", resourceAbsPath, err.Error())
		return
	}

	var profile string
	flag.StringVar(&profile, "gole.profile", "local", "环境变量")
	flag.Parse()
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		// 默认配置
		if profile == "" {
			fileName := file.Name()
			switch fileName {
			case "application.yaml":
				{
					LoadYamlFile(resourceAbsPath + "application.yaml")
					return
				}
			case "application.yml":
				{
					LoadYamlFile(resourceAbsPath + "application.yml")
					return
				}
			case "application.properties":
				{
					LoadYamlFile(resourceAbsPath + "application.properties")
					return
				}
			case "application.json":
				{
					LoadYamlFile(resourceAbsPath + "application.json")
					return
				}
			}
		} else {
			fileName := file.Name()
			currentProfile := getProfileFromFileName(fileName)
			if currentProfile == profile {
				extend := getFileExtension(fileName)
				extend = strings.ToLower(extend)
				switch extend {
				case "yaml":
					{
						LoadYamlFile(resourceAbsPath + fileName)
						return
					}
				case "yml":
					{
						LoadYamlFile(resourceAbsPath + fileName)
						return
					}
				case "properties":
					{
						LoadPropertyFile(resourceAbsPath + fileName)
						return
					}
				case "json":
					{
						LoadJsonFile(resourceAbsPath + fileName)
						return
					}
				}
			}
		}
	}
}

func getProfileFromFileName(fileName string) string {
	if strings.HasPrefix(fileName, "application-") {
		words := strings.SplitN(fileName, ".", 2)
		appNames := words[0]

		appNameAndProfile := strings.SplitN(appNames, "-", 2)
		return appNameAndProfile[1]
	}
	return ""
}

func getFileExtension(fileName string) string {
	if strings.Contains(fileName, ".") {
		words := strings.SplitN(fileName, ".", 2)
		return words[1]
	}
	return ""
}

func LoadYamlFile(filePath string) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("fail to read file:", err)
	}

	if appProperty == nil {
		appProperty = &ApplicationProperty{}
	}

	property, err := yaml.YamlToProperties(string(content))
	valueMap, _ := yaml.PropertiesToMap(property)
	appProperty.ValueMap = valueMap

	yamlMap, err := yaml.YamlToMap(string(content))
	appProperty.ValueDeepMap = yamlMap
}

func LoadPropertyFile(filePath string) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("fail to read file:", err)
	}

	if appProperty == nil {
		appProperty = &ApplicationProperty{}
	}

	valueMap, _ := yaml.PropertiesToMap(string(content))
	appProperty.ValueMap = valueMap

	yamlStr, _ := yaml.PropertiesToYaml(string(content))
	yamlMap, _ := yaml.YamlToMap(yamlStr)
	appProperty.ValueDeepMap = yamlMap
}

func LoadJsonFile(filePath string) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("fail to read file:", err)
	}

	if appProperty == nil {
		appProperty = &ApplicationProperty{}
	}

	yamlStr, err := yaml.JsonToYaml(string(content))
	property, err := yaml.YamlToProperties(yamlStr)
	valueMap, _ := yaml.PropertiesToMap(property)
	appProperty.ValueMap = valueMap

	yamlMap, _ := yaml.YamlToMap(yamlStr)
	appProperty.ValueDeepMap = yamlMap
}

func SetValue(key string, value interface{}) {
	if appProperty == nil {
		appProperty = &ApplicationProperty{}
		appProperty.ValueMap = map[string]interface{}{}
	}
	appProperty.ValueMap[key] = value
}

func GetValueString(key string) string {
	if value, exist := appProperty.ValueMap[key]; exist {
		return util.ToString(value)
	}
	return ""
}

func GetValueInt(key string) int {
	if value, exist := appProperty.ValueMap[key]; exist {
		return util.ToInt(value)
	}
	return 0
}

func GetValueInt8(key string) int8 {
	if value, exist := appProperty.ValueMap[key]; exist {
		return util.ToInt8(value)
	}
	return 0
}

func GetValueInt16(key string) int16 {
	if value, exist := appProperty.ValueMap[key]; exist {
		return util.ToInt16(value)
	}
	return 0
}

func GetValueInt32(key string) int32 {
	if value, exist := appProperty.ValueMap[key]; exist {
		return util.ToInt32(value)
	}
	return 0
}

func GetValueInt64(key string) int64 {
	if value, exist := appProperty.ValueMap[key]; exist {
		return util.ToInt64(value)
	}
	return 0
}

func GetValueUInt(key string) uint {
	if value, exist := appProperty.ValueMap[key]; exist {
		return util.ToUInt(value)
	}
	return 0
}

func GetValueUInt8(key string) uint8 {
	if value, exist := appProperty.ValueMap[key]; exist {
		return util.ToUInt8(value)
	}
	return 0
}

func GetValueUInt16(key string) uint16 {
	if value, exist := appProperty.ValueMap[key]; exist {
		return util.ToUInt16(value)
	}
	return 0
}

func GetValueUInt32(key string) uint32 {
	if value, exist := appProperty.ValueMap[key]; exist {
		return util.ToUInt32(value)
	}
	return 0
}

func GetValueUInt64(key string) uint64 {
	if value, exist := appProperty.ValueMap[key]; exist {
		return util.ToUInt64(value)
	}
	return 0
}

func GetValueFloat32(key string) float32 {
	if value, exist := appProperty.ValueMap[key]; exist {
		return util.ToFloat32(value)
	}
	return 0
}

func GetValueFloat64(key string) float64 {
	if value, exist := appProperty.ValueMap[key]; exist {
		return util.ToFloat64(value)
	}
	return 0
}

func GetValueBool(key string) bool {
	if value, exist := appProperty.ValueMap[key]; exist {
		return util.ToBool(value)
	}
	return false
}

func GetValueStringDefault(key, defaultValue string) string {
	if value, exist := appProperty.ValueMap[key]; exist {
		return util.ToString(value)
	}
	return defaultValue
}

func GetValueIntDefault(key string, defaultValue int) int {
	if value, exist := appProperty.ValueMap[key]; exist {
		return util.ToInt(value)
	}
	return defaultValue
}

func GetValueInt8Default(key string, defaultValue int8) int8 {
	if value, exist := appProperty.ValueMap[key]; exist {
		return util.ToInt8(value)
	}
	return defaultValue
}

func GetValueInt16Default(key string, defaultValue int16) int16 {
	if value, exist := appProperty.ValueMap[key]; exist {
		return util.ToInt16(value)
	}
	return defaultValue
}

func GetValueInt32Default(key string, defaultValue int32) int32 {
	if value, exist := appProperty.ValueMap[key]; exist {
		return util.ToInt32(value)
	}
	return defaultValue
}

func GetValueInt64Default(key string, defaultValue int64) int64 {
	if value, exist := appProperty.ValueMap[key]; exist {
		return util.ToInt64(value)
	}
	return defaultValue
}

func GetValueUIntDefault(key string, defaultValue uint) uint {
	if value, exist := appProperty.ValueMap[key]; exist {
		return util.ToUInt(value)
	}
	return defaultValue
}

func GetValueUInt8Default(key string, defaultValue uint8) uint8 {
	if value, exist := appProperty.ValueMap[key]; exist {
		return util.ToUInt8(value)
	}
	return defaultValue
}

func GetValueUInt16Default(key string, defaultValue uint16) uint16 {
	if value, exist := appProperty.ValueMap[key]; exist {
		return util.ToUInt16(value)
	}
	return defaultValue
}

func GetValueUInt32Default(key string, defaultValue uint32) uint32 {
	if value, exist := appProperty.ValueMap[key]; exist {
		return util.ToUInt32(value)
	}
	return defaultValue
}

func GetValueUInt64Default(key string, defaultValue uint64) uint64 {
	if value, exist := appProperty.ValueMap[key]; exist {
		return util.ToUInt64(value)
	}
	return defaultValue
}

func GetValueFloat32Default(key string, defaultValue float32) float32 {
	if value, exist := appProperty.ValueMap[key]; exist {
		return util.ToFloat32(value)
	}
	return defaultValue
}

func GetValueFloat64Default(key string, defaultValue float64) float64 {
	if value, exist := appProperty.ValueMap[key]; exist {
		return util.ToFloat64(value)
	}
	return defaultValue
}

func GetValueBoolDefault(key string, defaultValue bool) bool {
	if value, exist := appProperty.ValueMap[key]; exist {
		return util.ToBool(value)
	}
	return false
}

func GetValueObject(key string, targetPtrObj interface{}) error {
	data := doGetValue(appProperty.ValueDeepMap, key)
	err := util.DataToObject(data, targetPtrObj)
	if err != nil {
		return err
	}
	return nil
}

func GetValue(key string) interface{} {
	return doGetValue(appProperty.ValueDeepMap, key)
}

func doGetValue(parentValue interface{}, key string) interface{} {
	if key == "" {
		return parentValue
	}
	parentValueKind := reflect.ValueOf(parentValue).Kind()
	if parentValueKind == reflect.Map {
		keys := strings.SplitN(key, ".", 2)
		v1 := reflect.ValueOf(parentValue).MapIndex(reflect.ValueOf(keys[0]))
		emptyValue := reflect.Value{}
		if v1 == emptyValue {
			return nil
		}
		if len(keys) == 1 {
			return doGetValue(v1.Interface(), "")
		} else {
			return doGetValue(v1.Interface(), fmt.Sprintf("%v", keys[1]))
		}
	}
	return nil
}

type ApplicationProperty struct {
	ValueMap     map[string]interface{}
	ValueDeepMap map[string]interface{}
}
