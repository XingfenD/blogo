package tpl

import (
	"text/template"

	"github.com/XingfenD/blogo/module/loader"
)

func loadColleTpl() {
	var err error
	ColleTpl = template.New("collection.html").Funcs(funcMap)
	ColleTpl, err = ColleTpl.ParseFiles(
		rootPath+"/template/page/admin/collection.html",
		rootPath+"/template/layout/footer.html",
		rootPath+"/template/layout/sidebar.html",
	)
	if err != nil || ColleTpl == nil {
		loader.Logger.Error("load admin/collection template failed", err)
	} else {
		loader.Logger.Info("load admin/collection template success")
	}
}

func loadPostTableTpl() {
	var err error
	PostTableTpl = template.New("post-collection.html").Funcs(funcMap)
	PostTableTpl, err = PostTableTpl.ParseFiles(
		rootPath+"/template/page/admin/post-collection.html",
		rootPath+"/template/layout/footer.html",
		rootPath+"/template/layout/sidebar.html",
	)
	if err != nil || PostTableTpl == nil {
		loader.Logger.Error("load admin/post-table template failed", err)
	} else {
		loader.Logger.Info("load admin/post-table template success")
	}
}
