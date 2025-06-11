package sqlite

import (
	"database/sql"

	"github.com/XingfenD/blogo/module/loader"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB(db_path string) error {
	var err error
	loader.Logger.Info("Initializing database...")
	db, err = sql.Open("sqlite3", db_path)
	if err != nil {
		loader.Logger.Error("Error opening database:", err)
		return err
	}
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS "blog_posts" (
			"blog_id" INTEGER NOT NULL UNIQUE,
			"dir_name" TEXT NOT NULL UNIQUE,
			"title" TEXT,
			"description" TEXT,
			"content" TEXT NOT NULL,
			"cate_id" INTEGER,
			"create_time" TEXT NOT NULL,
			"last_modified" TEXT NOT NULL,
			PRIMARY KEY("blog_id"),
			FOREIGN KEY ("cate_id") REFERENCES "categories"("cate_id")
			ON UPDATE NO ACTION ON DELETE NO ACTION
		);

		CREATE TABLE IF NOT EXISTS "tags" (
			"id" INTEGER NOT NULL UNIQUE,
			"name" TEXT,
			PRIMARY KEY("id")
		);

		CREATE TABLE IF NOT EXISTS "blog_tag" (
			"blog_id" INTEGER NOT NULL,
			"tag_id" INTEGER NOT NULL,
			PRIMARY KEY("blog_id", "tag_id"),
			FOREIGN KEY ("tag_id") REFERENCES "tags"("id")
			ON UPDATE NO ACTION ON DELETE NO ACTION,
			FOREIGN KEY ("blog_id") REFERENCES "blog_posts"("blog_id")
			ON UPDATE NO ACTION ON DELETE NO ACTION
		);

		CREATE TABLE IF NOT EXISTS "categories" (
			"cate_id" INTEGER NOT NULL UNIQUE,
			"cate_name" TEXT UNIQUE,
			PRIMARY KEY("cate_id")
		);
	`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		loader.Logger.Error("Error creating table:", err)
		return err
	}
	loader.Logger.Info("Database initialized successfully")
	// return exampleQuery()
	return nil
}

// func exampleQuery() error {
// 	exampleSQL := `
// 		INSERT INTO "categories" ("cate_id", "cate_name") VALUES
// 		(1, 'Technology'),
// 		(2, 'Lifestyle'),
// 		(3, 'Travel'),
// 		(4, 'Food');

// 		INSERT INTO "tags" ("id", "name") VALUES
// 		(1, 'Programming'),
// 		(2, 'Health'),
// 		(3, 'Nature'),
// 		(4, 'Cooking'),
// 		(5, 'History'),
// 		(6, 'Adventure');

// 		INSERT INTO "blog_posts" ("blog_id", "dir_name", "title", "description", "content", "cate_id", "create_time", "last_modified") VALUES
// 		(1, 'about', '关于blogo', 'blogo简介', 'blogo内容', NULL, '2025-06-01 10:00:00', '2025-06-01 10:00:00'),
// 		(2, 'LifeHacks', 'Healthy Living Tips', 'Tips for a healthy lifestyle', 'Eating well and exercising regularly are key...', 2, '2025-06-02 11:00:00', '2025-06-02 11:00:00'),
// 		(3, 'TravelDiary', 'Exploring Europe', 'A journey through Europe', 'Visiting historic cities and beautiful landscapes...', 3, '2025-06-03 12:00:00', '2025-06-03 12:00:00'),
// 		(4, 'FoodieBlog', 'Best Italian Recipes', 'Delicious Italian dishes', 'From pasta to pizza, Italian cuisine is amazing...', 4, '2025-06-04 13:00:00', '2025-06-04 13:00:00');

// 		INSERT INTO "blog_tag" ("blog_id", "tag_id") VALUES
// 		(1, 1), -- TechBlog tagged with Programming
// 		(2, 2), -- LifeHacks tagged with Health
// 		(3, 3), -- TravelDiary tagged with Nature
// 		(4, 4), -- FoodieBlog tagged with Cooking
// 		(1, 5), -- TechBlog also tagged with History
// 		(3, 6); -- TravelDiary also tagged with Adventure
// 	`
// 	_, err := db.Exec(exampleSQL)
// 	if err != nil {
// 		loader.Logger.Error("Error executing example query:", err)
// 		return err
// 	}
// 	return nil
// }

func CloseDB() {
	err := db.Close()
	if err != nil {
		loader.Logger.Error("Error closing database:", err)
	}
}

type ArticleMeta struct {
	Title        string
	CreateDate   string
	LastModified string
	Tags         []string
	Category     struct {
		Name string
		Id   int
	}
	Description string
	Content     string
}

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
        SELECT t.name
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
		if err := tagsRows.Scan(&tag); err != nil {
			continue
		}
		meta.Tags = append(meta.Tags, tag)
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
        SELECT t.name
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
		if err := tagsRows.Scan(&tag); err != nil {
			continue
		}
		meta.Tags = append(meta.Tags, tag)
	}

	return &meta, nil
}

func GetCategoryList() []struct {
	Name string
	Id   int
	Time string
} {
	var categories []struct {
		Name string
		Id   int
		Time string
	}

	rows, err := db.Query(`
        SELECT cate_id, cate_name
        FROM categories
        ORDER BY cate_id`)
	if err != nil {
		loader.Logger.Error("Error querying categories:", err)
		return categories
	}
	defer rows.Close()

	for rows.Next() {
		var cate struct {
			Name string
			Id   int
		}
		if err := rows.Scan(&cate.Id, &cate.Name); err != nil {
			loader.Logger.Error("Error scanning category row:", err)
			continue
		}
		categories = append(categories, struct {
			Name string
			Id   int
			Time string
		}{
			Name: cate.Name,
			Id:   cate.Id,
			Time: "",
		})
	}

	if err = rows.Err(); err != nil {
		loader.Logger.Error("Error after scanning categories:", err)
	}

	return categories
}

func GetCateById(cateId int) (string, error) {
	var cateName string
	err := db.QueryRow(`
        SELECT cate_name
        FROM categories
        WHERE cate_id =?`, cateId).Scan(&cateName)
	if err != nil {
		loader.Logger.Error("Error querying category by ID:", err)
		return "", err
	}
	return cateName, nil
}

func GetTagList() []struct {
	Name string
	Id   int
} {
	var tags []struct {
		Name string
		Id   int
	}

	rows, err := db.Query(`
        SELECT id, name
        FROM tags
        ORDER BY id`)
	if err != nil {
		loader.Logger.Error("Error querying tags:", err)
		return tags
	}
	defer rows.Close()

	for rows.Next() {
		var tag struct {
			Name string
			Id   int
		}
		if err := rows.Scan(&tag.Id, &tag.Name); err != nil {
			loader.Logger.Error("Error scanning tag row:", err)
			continue
		}
		tags = append(tags, tag)
	}

	if err = rows.Err(); err != nil {
		loader.Logger.Error("Error after scanning tags:", err)
	}

	return tags
}

func GetTagById(tagId int) (string, error) {
	var tagName string
	err := db.QueryRow(`
        SELECT name
        FROM tags
        WHERE id =?`, tagId).Scan(&tagName)
	if err != nil {
		loader.Logger.Error("Error querying tag by ID:", err)
		return "", err
	}
	return tagName, nil
}

func GetArticlesByCategory(cateId int) []struct {
	BlogId       int
	DirName      string
	Title        string
	CreateDate   string
	LastModified string
	Tags         []string
	Description  string
	Content      string
} {
	var articles []struct {
		BlogId       int
		DirName      string
		Title        string
		CreateDate   string
		LastModified string
		Tags         []string
		Description  string
		Content      string
	}

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
		var article struct {
			BlogId       int
			DirName      string
			Title        string
			CreateDate   string
			LastModified string
			Tags         []string
			Description  string
			Content      string
		}

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
		articles = append(articles, struct {
			BlogId       int
			DirName      string
			Title        string
			CreateDate   string
			LastModified string
			Tags         []string
			Description  string
			Content      string
		}{
			BlogId:       article.BlogId,
			DirName:      article.DirName, // 新增字段赋值
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

func GetArticlesByTag(tagId int) []struct {
	BlogId       int
	DirName      string
	Title        string
	CreateDate   string
	LastModified string
	Tags         []string
	Description  string
	Content      string
} {
	var articles []struct {
		BlogId       int
		DirName      string
		Title        string
		CreateDate   string
		LastModified string
		Tags         []string
		Description  string
		Content      string
	}

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
		var article struct {
			BlogId       int
			DirName      string
			Title        string
			CreateDate   string
			LastModified string
			Tags         []string
			Description  string
			Content      string
		}

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

func GetArticleList() []struct {
	BlogId       int
	DirName      string
	Title        string
	CreateDate   string
	LastModified string
	Year         string // 新增年份字段
	Tags         []string
	Description  string
	Content      string
} {
	var articles []struct {
		BlogId       int
		DirName      string
		Title        string
		CreateDate   string
		LastModified string
		Year         string
		Tags         []string
		Description  string
		Content      string
	}

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
		var article struct {
			BlogId       int
			DirName      string
			Title        string
			CreateDate   string
			LastModified string
			Year         string
			Tags         []string
			Description  string
			Content      string
		}

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
		// ... 保持后续标签查询逻辑不变 ...
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
