package sqlite_db

import (
	"database/sql"

	"github.com/XingfenD/blogo/module/loader"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

type ArticleMeta struct {
	Title        string
	CreateDate   string
	LastModified string
	Tags         []struct {
		Name string
		Url  string
	}
	Category struct {
		Name string
		Id   int
	}
	Description string
	Content     string
}

type ArticleListItem struct {
	BlogId       int
	DirName      string
	Title        string
	CreateDate   string
	LastModified string
	Year         string
	Tags         []string
	Category     struct {
		Name string
		Id   int
	}
	Description string
	Content     string
}

type CollectionListItem struct {
	ColleId   int
	ColleName string
	PostCount int
}

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

func CloseDB() {
	err := db.Close()
	if err != nil {
		loader.Logger.Error("Error closing database:", err)
	}
}
