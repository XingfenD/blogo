package tpl

import (
	"text/template"

	"github.com/XingfenD/blogo/module/loader"
)

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
