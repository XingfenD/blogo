package tpl

import (
	"text/template"

	"github.com/XingfenD/blogo/module/loader"
)

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
