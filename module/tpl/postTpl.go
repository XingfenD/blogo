package tpl

import (
	"text/template"

	"github.com/XingfenD/blogo/module/loader"
)

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
