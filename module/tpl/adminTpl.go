package tpl

import (
	"text/template"

	"github.com/XingfenD/blogo/module/loader"
)

func loadAdminTpl() {
	var err error
	AdminTpl = template.New("admin.html").Funcs(funcMap)
	AdminTpl, err = AdminTpl.ParseFiles(
		rootPath+"/template/page/admin/admin.html",
		rootPath+"/template/layout/dashboard.html",
		rootPath+"/template/layout/footer.html",
		rootPath+"/template/layout/sidebar.html",
	)
	if err != nil || AdminTpl == nil {
		loader.Logger.Error("load admin/admin template failed", err)
	} else {
		loader.Logger.Info("load admin/admin template success")
	}
}
