package tpl

import (
	"text/template"

	"github.com/XingfenD/blogo/module/loader"
)

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
