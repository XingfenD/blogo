package sqlite_db

import (
	"database/sql"
	"fmt"

	"github.com/XingfenD/blogo/module/loader"
)

func GetAboutMeta() (*ArticleMeta, error) {
	// 查询基础文章信息
	row := db.QueryRow(`
        SELECT title, description, content, create_time, last_modified, cate_id
        FROM blog_posts
        WHERE dir_name = 'about'`)

	var meta ArticleMeta
	var cateID sql.NullInt64 // 处理可能为NULL的分类ID

	err := row.Scan(
		&meta.Title,
		&meta.Description,
		&meta.Content,
		&meta.CreateDate,
		&meta.LastModified,
		&cateID,
	)
	if err != nil {
		loader.Logger.Error("Error querying about post:", err)
		return nil, err
	}

	// 查询分类信息（如果存在）
	if cateID.Valid {
		err = db.QueryRow(`
            SELECT cate_name
            FROM categories
            WHERE cate_id = ?`, cateID.Int64).Scan(&meta.Category.Name)
		if err != nil {
			loader.Logger.Error("Error querying category:", err)
		}
		meta.Category.Id = int(cateID.Int64)
	}

	// 查询标签信息
	tagsRows, err := db.Query(`
        SELECT t.name, t.id
        FROM tags t
        INNER JOIN blog_tag bt ON t.id = bt.tag_id
        INNER JOIN blog_posts bp ON bt.blog_id = bp.blog_id
        WHERE bp.dir_name = 'about'`)
	if err != nil {
		loader.Logger.Error("Error querying tags:", err)
		return &meta, nil
	}
	defer tagsRows.Close()

	for tagsRows.Next() {
		var tag string
		var tagId int
		if err := tagsRows.Scan(&tag, &tagId); err != nil {
			continue
		}
		meta.Tags = append(meta.Tags, struct {
			Name string
			Url  string
		}{
			Name: tag,
			Url:  fmt.Sprintf("archives/tags/%d", tagId),
		})
	}

	return &meta, nil
}

func GetArticleMetaByDir(dirName string) (*ArticleMeta, error) {
	// 查询基础文章信息
	row := db.QueryRow(`
        SELECT title, description, content, create_time, last_modified, cate_id
        FROM blog_posts
        WHERE dir_name = ?`, dirName)

	var meta ArticleMeta
	var cateID sql.NullInt64

	err := row.Scan(
		&meta.Title,
		&meta.Description,
		&meta.Content,
		&meta.CreateDate,
		&meta.LastModified,
		&cateID,
	)
	if err != nil {
		loader.Logger.Error("Error querying article by dir:", err)
		return nil, err
	}

	// 查询分类信息（如果存在）
	if cateID.Valid {
		err = db.QueryRow(`
            SELECT cate_name
            FROM categories
            WHERE cate_id = ?`, cateID.Int64).Scan(&meta.Category.Name)
		if err != nil {
			loader.Logger.Error("Error querying category:", err)
		}
		meta.Category.Id = int(cateID.Int64)
	}

	// 查询标签信息
	tagsRows, err := db.Query(`
        SELECT t.name, t.id
        FROM tags t
        INNER JOIN blog_tag bt ON t.id = bt.tag_id
        INNER JOIN blog_posts bp ON bt.blog_id = bp.blog_id
        WHERE bp.dir_name = ?`, dirName)
	if err != nil {
		loader.Logger.Error("Error querying tags:", err)
		return &meta, nil
	}
	defer tagsRows.Close()

	for tagsRows.Next() {
		var tag string
		var tagId int
		if err := tagsRows.Scan(&tag, &tagId); err != nil {
			continue
		}
		meta.Tags = append(meta.Tags, struct {
			Name string
			Url  string
		}{
			Name: tag,
			Url:  fmt.Sprintf("archives/tags/%d", tagId),
		})
	}

	return &meta, nil
}

func GetArticlesByCategory(cateId int) []ArticleListItem {
	var articles []ArticleListItem

	// 查询基础文章信息

	rows, err := db.Query(`
        SELECT b.title, b.description, b.content, b.create_time, b.last_modified,
               b.blog_id, b.dir_name
        FROM blog_posts b
        WHERE b.cate_id = ?
        ORDER BY b.create_time DESC`, cateId)
	if err != nil {
		loader.Logger.Error("Error querying articles:", err)
		return articles
	}
	defer rows.Close()

	for rows.Next() {
		var article ArticleListItem

		err := rows.Scan(
			&article.Title,
			&article.Description,
			&article.Content,
			&article.CreateDate,
			&article.LastModified,
			&article.BlogId,
			&article.DirName,
		)
		if err != nil {
			loader.Logger.Error("Error scanning article row:", err)
			continue
		}

		// 查询标签信息
		tagRows, err := db.Query(`
            SELECT t.name
            FROM tags t
            INNER JOIN blog_tag bt ON t.id = bt.tag_id
            WHERE bt.blog_id = ?`, article.BlogId)
		if err != nil {
			loader.Logger.Error("Error querying tags:", err)
			continue
		}

		for tagRows.Next() {
			var tag string
			if err := tagRows.Scan(&tag); err == nil {
				article.Tags = append(article.Tags, tag)
			}
		}
		tagRows.Close()

		// 修改返回结构
		articles = append(articles, ArticleListItem{
			BlogId:       article.BlogId,
			DirName:      article.DirName,
			Title:        article.Title,
			CreateDate:   article.CreateDate,
			LastModified: article.LastModified,
			Tags:         article.Tags,
			Description:  article.Description,
			Content:      article.Content,
		})
	}

	if err = rows.Err(); err != nil {
		loader.Logger.Error("Error after scanning articles:", err)
	}

	return articles
}

func GetArticlesByTag(tagId int) []ArticleListItem {
	var articles []ArticleListItem

	// 查询基础文章信息
	rows, err := db.Query(`
        SELECT b.title, b.description, b.content, b.create_time, b.last_modified,
               b.blog_id, b.dir_name
        FROM blog_posts b
        INNER JOIN blog_tag bt ON b.blog_id = bt.blog_id
        WHERE bt.tag_id = ?
        ORDER BY b.create_time DESC`, tagId)
	if err != nil {
		loader.Logger.Error("Error querying articles:", err)
		return articles
	}
	defer rows.Close()

	for rows.Next() {
		var article ArticleListItem

		err := rows.Scan(
			&article.Title,
			&article.Description,
			&article.Content,
			&article.CreateDate,
			&article.LastModified,
			&article.BlogId,
			&article.DirName,
		)
		if err != nil {
			loader.Logger.Error("Error scanning article row:", err)
			continue
		}

		// 查询标签信息
		tagRows, err := db.Query(`
            SELECT t.name
            FROM tags t
            INNER JOIN blog_tag bt ON t.id = bt.tag_id
            WHERE bt.blog_id = ?`, article.BlogId)
		if err != nil {
			loader.Logger.Error("Error querying tags:", err)
			continue
		}

		for tagRows.Next() {
			var tag string
			if err := tagRows.Scan(&tag); err == nil {
				article.Tags = append(article.Tags, tag)
			}
		}
		tagRows.Close()

		articles = append(articles, article)
	}

	if err = rows.Err(); err != nil {
		loader.Logger.Error("Error after scanning articles:", err)
	}

	return articles
}

func GetArticleList() []ArticleListItem {
	var articles []ArticleListItem

	// 修改查询语句添加年份字段和排序规则
	rows, err := db.Query(`
        SELECT
            b.title,
            b.description,
            b.content,
            b.create_time,
            b.last_modified,
            b.blog_id,
            b.dir_name,
            strftime('%Y', create_time) AS year
        FROM blog_posts b
        ORDER BY year DESC, b.create_time DESC`) // 按年份降序+时间降序
	if err != nil {
		loader.Logger.Error("Error querying articles:", err)
		return articles
	}
	defer rows.Close()

	for rows.Next() {
		var article ArticleListItem

		err := rows.Scan(
			&article.Title,
			&article.Description,
			&article.Content,
			&article.CreateDate,
			&article.LastModified,
			&article.BlogId,
			&article.DirName,
			&article.Year, // 扫描新增的年份字段
		)

		if err != nil {
			loader.Logger.Error("Error scanning article row:", err)
			continue
		}

		// 查询标签信息（与 GetArticlesByCategory 相同逻辑）
		tagRows, err := db.Query(`
            SELECT t.name
            FROM tags t
            INNER JOIN blog_tag bt ON t.id = bt.tag_id
            WHERE bt.blog_id = ?`, article.BlogId)
		if err != nil {
			loader.Logger.Error("Error querying tags:", err)
			continue
		}

		for tagRows.Next() {
			var tag string
			if err := tagRows.Scan(&tag); err == nil {
				article.Tags = append(article.Tags, tag)
			}
		}
		tagRows.Close()

		articles = append(articles, article)
	}

	if err = rows.Err(); err != nil {
		loader.Logger.Error("Error after scanning articles:", err)
	}

	return articles
}

func GetRecentArticles(limit int) []ArticleListItem {
	var articles []ArticleListItem

	rows, err := db.Query(`
        SELECT
            b.title,
            b.dir_name,
            b.create_time,
            b.last_modified,
            c.cate_name
        FROM blog_posts b
        LEFT JOIN categories c ON b.cate_id = c.cate_id
        ORDER BY b.last_modified DESC
        LIMIT ?`, limit)
	if err != nil {
		loader.Logger.Error("Error querying recent articles:", err)
		return articles
	}
	defer rows.Close()

	for rows.Next() {
		var article ArticleListItem
		var category sql.NullString

		err := rows.Scan(
			&article.Title,
			&article.DirName,
			&article.CreateDate,
			&article.LastModified,
			&category,
		)

		if err == nil {
			if category.Valid {
				article.Category.Name = category.String
			}
			articles = append(articles, article)
		}
	}

	if err = rows.Err(); err != nil {
		loader.Logger.Error("Error after scanning recent articles:", err)
	}

	return articles
}
