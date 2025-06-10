package main

import (
	"database/sql"

	"github.com/XingfenD/blogo/module/loader"
	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func InitDB() error {
	var err error
	db, err = sql.Open("sqlite3", "./blogo_db.db")
	if err != nil {
		loader.Logger.Error("Error opening database:", err)
		return err
	}
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS "blog_posts" (
			"blog_id" INTEGER NOT NULL UNIQUE,
			"title" TEXT NOT NULL,
			"content" TEXT,
			"cate_id" INTEGER,
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
			"blog_id" INTEGER,
			"tag_id" INTEGER,
			PRIMARY KEY("blog_id", "tag_id"),
			FOREIGN KEY ("tag_id") REFERENCES "tags"("id")
			ON UPDATE NO ACTION ON DELETE NO ACTION,
			FOREIGN KEY ("blog_id") REFERENCES "blog_posts"("blog_id")
			ON UPDATE NO ACTION ON DELETE NO ACTION
		);

		CREATE TABLE IF NOT EXISTS "categories" (
			"cate_id" INTEGER NOT NULL UNIQUE,
			"cate_name" TEXT,
			PRIMARY KEY("cate_id")
		);
	`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		loader.Logger.Error("Error creating table:", err)
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

// func create_category(cate_name string) error {

// }

// func update_category(cate_id int, cate_name string) error {

// }

// func delete_category(cate_id int) error {

// }

// func select_category(cate_id int) error {

// }
