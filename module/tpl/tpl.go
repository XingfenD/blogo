package tpl

import (
	"reflect"
	"sort"
	"strconv"
	"text/template"
	"time"

	"github.com/XingfenD/blogo/module/loader"
)

var IndexTpl *template.Template
var PostTpl *template.Template
var SectionTpl *template.Template
var ArchiveTpl *template.Template
var NotFoundTpl *template.Template
var AdminTpl *template.Template
var ColleTpl *template.Template
var PostTableTpl *template.Template

var rootPath string

var funcMap = template.FuncMap{
	"date": func(format string) string {
		return time.Now().Format(format)
	},
	"sortByYearDesc": func(keys []string) []string {
		// 将字符串年份转换为数字排序
		years := make([]int, len(keys))
		for i, k := range keys {
			y, _ := strconv.Atoi(k)
			years[i] = y
		}

		sort.Sort(sort.Reverse(sort.IntSlice(years)))

		// 转换回字符串
		result := make([]string, len(years))
		for i, y := range years {
			result[i] = strconv.Itoa(y)
		}
		return result
	},
	"getMapValue": func(m interface{}, key string) interface{} {
		val := reflect.ValueOf(m)
		if val.Kind() != reflect.Map {
			return nil
		}

		mapValue := val.MapIndex(reflect.ValueOf(key))
		if mapValue.IsValid() {
			return mapValue.Interface()
		}
		return nil
	},
	"keys": func(m interface{}) []string {
		val := reflect.ValueOf(m)
		if val.Kind() != reflect.Map {
			return nil
		}

		keys := make([]string, 0, val.Len())
		for _, k := range val.MapKeys() {
			if k.Kind() == reflect.String {
				keys = append(keys, k.String())
			}
		}
		return keys
	},
}

func LoadTemplate(root_path string) {
	loader.Logger.Info("loading templates")
	rootPath = root_path

	loadIndexTpl()
	loadPostTpl()
	loadSectionTpl()
	loadArchiveTpl()
	loadAdminTpl()
	load404Tpl()
	loadColleTpl()
	loadPostTableTpl()

	loader.Logger.Info("load templates success")
}
