package tpl

import (
	"text/template"

	"github.com/XingfenD/blogo/module/loader"
)

func loadColleTpl() {
	var err error
	ColleTpl = template.New("colle.html").Funcs(funcMap)
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
