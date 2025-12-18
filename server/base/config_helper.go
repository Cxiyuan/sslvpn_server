package base

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/spf13/viper"
)

var (
	cfgMu sync.RWMutex
)

// LoadConfigFromViper 从 viper 实例加载配置到 Cfg
// 支持的类型: string, int, bool
// 如需支持更多类型(如 time.Duration, []string 等)，请扩展此函数
func LoadConfigFromViper(v *viper.Viper) error {
	cfgMu.Lock()
	defer cfgMu.Unlock()

	// 使用反射更新 Cfg 的值
	cfgValue := reflect.ValueOf(Cfg).Elem()
	cfgType := cfgValue.Type()

	for i := 0; i < cfgValue.NumField(); i++ {
		field := cfgValue.Field(i)
		fieldType := cfgType.Field(i)
		
		// 获取 json tag 作为配置键名
		tag := fieldType.Tag.Get("json")
		if tag == "" || tag == "-" {
			continue
		}

		// 跳过不可设置的字段
		if !field.CanSet() {
			continue
		}

		// 根据字段类型设置值
		// 注意：当前仅支持 string, int, bool 三种基础类型
		// 如需支持更多类型，请在下方 switch 中添加对应 case
		switch field.Kind() {
		case reflect.String:
			if v.IsSet(tag) {
				field.SetString(v.GetString(tag))
			}
		case reflect.Int:
			if v.IsSet(tag) {
				field.SetInt(int64(v.GetInt(tag)))
			}
		case reflect.Bool:
			if v.IsSet(tag) {
				field.SetBool(v.GetBool(tag))
			}
		default:
			// 不支持的类型将被跳过，不会报错
			// 可在此处添加日志记录
		}
	}

	return nil
}

// GetConfigValue 安全读取配置值
func GetConfigValue(getter func(*ServerConfig) interface{}) interface{} {
	cfgMu.RLock()
	defer cfgMu.RUnlock()
	return getter(Cfg)
}

// ServerCfg2Map 将服务器配置转换为前端所需的格式
func ServerCfg2Map() []map[string]interface{} {
	var result []map[string]interface{}

	for _, v := range configs {
		var val interface{}
		
		// 获取当前配置值
		cfgValue := reflect.ValueOf(Cfg).Elem()
		cfgType := cfgValue.Type()
		
		for i := 0; i < cfgValue.NumField(); i++ {
			fieldType := cfgType.Field(i)
			tag := fieldType.Tag.Get("json")
			
			if tag == v.Name {
				field := cfgValue.Field(i)
				val = field.Interface()
				break
			}
		}

		// 如果没有从 Cfg 中找到值，使用默认值
		if val == nil {
			switch v.Typ {
			case cfgStr:
				val = v.ValStr
			case cfgInt:
				val = v.ValInt
			case cfgBool:
				val = v.ValBool
			}
		}

		item := map[string]interface{}{
			"name":      v.Name,
			"env":       fmt.Sprintf("LINK_%s", envs[v.Name]),
			"info":      v.Usage,
			"val":       val,
			"sensitive": v.Sensitive, // 添加敏感标记
		}

		result = append(result, item)
	}

	return result
}