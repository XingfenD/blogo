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
			"blog_name" TEXT NOT NULL UNIQUE,
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

func exampleQuery() error {
	exampleSQL := `
		INSERT INTO "categories" ("cate_id", "cate_name") VALUES
		(1, 'Technology'),
		(2, 'Lifestyle'),
		(3, 'Travel'),
		(4, 'Food');

		INSERT INTO "tags" ("id", "name") VALUES
		(1, 'Programming'),
		(2, 'Health'),
		(3, 'Nature'),
		(4, 'Cooking'),
		(5, 'History'),
		(6, 'Adventure');

		INSERT INTO "blog_posts" ("blog_id", "blog_name", "title", "description", "content", "cate_id", "create_time", "last_modified") VALUES
		(1, 'about', '关于blogo', 'blogo简介', 'blogo内容', NULL, '2025-06-01 10:00:00', '2025-06-01 10:00:00'),
		(2, 'LifeHacks', 'Healthy Living Tips', 'Tips for a healthy lifestyle', 'Eating well and exercising regularly are key...', 2, '2025-06-02 11:00:00', '2025-06-02 11:00:00'),
		(3, 'TravelDiary', 'Exploring Europe', 'A journey through Europe', 'Visiting historic cities and beautiful landscapes...', 3, '2025-06-03 12:00:00', '2025-06-03 12:00:00'),
		(4, 'FoodieBlog', 'Best Italian Recipes', 'Delicious Italian dishes', 'From pasta to pizza, Italian cuisine is amazing...', 4, '2025-06-04 13:00:00', '2025-06-04 13:00:00');

		INSERT INTO "blog_tag" ("blog_id", "tag_id") VALUES
		(1, 1), -- TechBlog tagged with Programming
		(2, 2), -- LifeHacks tagged with Health
		(3, 3), -- TravelDiary tagged with Nature
		(4, 4), -- FoodieBlog tagged with Cooking
		(1, 5), -- TechBlog also tagged with History
		(3, 6); -- TravelDiary also tagged with Adventure
	`
	_, err := db.Exec(exampleSQL)
	if err != nil {
		loader.Logger.Error("Error executing example query:", err)
		return err
	}
	return nil
}

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
        WHERE blog_name = 'about'`)

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
        WHERE bp.blog_name = 'about'`)
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
