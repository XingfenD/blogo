package tpl

import (
	"text/template"

	"github.com/XingfenD/blogo/module/loader"
)

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
