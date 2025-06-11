package tpl

import (
	"text/template"
	"time"

	"github.com/XingfenD/blogo/module/loader"
)

var IndexTpl *template.Template
var PostTpl *template.Template
var SectionTpl *template.Template
var ArchiveTpl *template.Template

var rootPath string

var funcMap = template.FuncMap{
	"date": func(format string) string {
		return time.Now().Format(format)
	},
}

func LoadTemplate(root_path string) {
	loader.Logger.Info("loading template")
	rootPath = root_path

	loadIndexTpl()
	loadPostTpl()
	loadSectionTpl()
	loadArchiveTpl()

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
	if err != nil {
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
	if err != nil {
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
	if err != nil {
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
	if err != nil {
		loader.Logger.Error("load index template failed", err)
	} else {
		loader.Logger.Info("load index template success")
	}
}
