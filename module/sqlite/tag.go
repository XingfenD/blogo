package sqlite_db

import (
	"github.com/XingfenD/blogo/module/loader"
)

func GetTagList(needCount bool) []CollectionListItem {
	var tags []CollectionListItem

	// 基础查询语句
	query := `
        SELECT id, name
        FROM tags
        ORDER BY id`

	if needCount {
		query = `
            SELECT tags.id, tags.name, COUNT(blog_tag.blog_id) AS post_count
            FROM tags
            LEFT JOIN blog_tag ON tags.id = blog_tag.tag_id
            GROUP BY tags.id
            ORDER BY tags.id`
	}

	rows, err := db.Query(query)
	if err != nil {
		loader.Logger.Error("Error querying tags:", err)
		return tags
	}
	defer rows.Close()

	for rows.Next() {
		var tag CollectionListItem
		if needCount {
			err = rows.Scan(&tag.ColleId, &tag.ColleName, &tag.PostCount)
		} else {
			err = rows.Scan(&tag.ColleId, &tag.ColleName)
		}
		if err != nil {
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
