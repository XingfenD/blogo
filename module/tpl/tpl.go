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
	loader.Logger.Info("loading template")
	rootPath = root_path

	loadIndexTpl()
	loadPostTpl()
	loadSectionTpl()
	loadArchiveTpl()
	load404Tpl()

	loader.Logger.Info("load template success")
}

func loadPostTpl() {
	var err error
	PostTpl = template.New("post.html").Funcs(funcMap)
	PostTpl, err = PostTpl.ParseFiles(
		rootPath+"/template/page/post.html",
		rootPath+"/template/layout/footer.html",
		rootPath+"/template/layout/sidebar.html",
		rootPath+"/template/layout/article.html",
	)
	if err != nil || PostTpl == nil {
		loader.Logger.Error("load post template failed", err)
	} else {
		loader.Logger.Info("load post template success")
	}

}

func loadSectionTpl() {
	var err error
	SectionTpl = template.New("section.html").Funcs(funcMap)
	SectionTpl, err = SectionTpl.ParseFiles(
		rootPath+"/template/page/section.html",
		rootPath+"/template/layout/footer.html",
		rootPath+"/template/layout/sidebar.html",
	)
	if err != nil || SectionTpl == nil {
		loader.Logger.Error("load section template failed", err)
	} else {
		loader.Logger.Info("load section template success")
	}
}

func loadArchiveTpl() {
	var err error
	ArchiveTpl = template.New("archives.html").Funcs(funcMap)
	ArchiveTpl, err = ArchiveTpl.ParseFiles(
		rootPath+"/template/page/archives.html",
		rootPath+"/template/layout/footer.html",
		rootPath+"/template/layout/sidebar.html",
	)
	if err != nil || ArchiveTpl == nil {
		loader.Logger.Error("load archive template failed", err)
	} else {
		loader.Logger.Info("load archive template success")
	}
}

func loadIndexTpl() {
	var err error
	IndexTpl = template.New("index.html").Funcs(funcMap)
	IndexTpl, err = IndexTpl.ParseFiles(
		rootPath+"/template/page/index.html",
		rootPath+"/template/layout/footer.html",
		rootPath+"/template/layout/sidebar.html",
	)
	if err != nil || IndexTpl == nil {
		loader.Logger.Error("load index template failed", err)
	} else {
		loader.Logger.Info("load index template success")
	}
}

func load404Tpl() {
	var err error
	NotFoundTpl = template.New("404.html").Funcs(funcMap)
	NotFoundTpl, err = NotFoundTpl.ParseFiles(
		rootPath+"/template/page/404.html",
		rootPath+"/template/layout/footer.html",
		rootPath+"/template/layout/sidebar.html",
	)
	if err != nil || NotFoundTpl == nil {
		loader.Logger.Error("load 404 template failed", err)
	} else {
		loader.Logger.Info("load 404 template success")
	}
}
