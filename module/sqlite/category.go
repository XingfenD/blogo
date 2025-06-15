package sqlite_db

import "github.com/XingfenD/blogo/module/loader"

func GetCategoryList(needCount bool) []CollectionListItem {
	var categories []CollectionListItem

	// 动态构建查询语句
	query := `
        SELECT c.cate_id, c.cate_name
        FROM categories c
        ORDER BY c.cate_id`
	if needCount {
		query = `
            SELECT c.cate_id, c.cate_name, COUNT(b.blog_id) AS post_count
            FROM categories c
            LEFT JOIN blog_posts b ON c.cate_id = b.cate_id
            GROUP BY c.cate_id
            ORDER BY c.cate_id`
	}

	rows, err := db.Query(query)
	if err != nil {
		loader.Logger.Error("Error querying categories:", err)
		return categories
	}
	defer rows.Close()

	for rows.Next() {
		var (
			id        int
			name      string
			postCount int
		)

		if needCount {
			err = rows.Scan(&id, &name, &postCount)
		} else {
			err = rows.Scan(&id, &name)
		}

		if err != nil {
			loader.Logger.Error("Error scanning category row:", err)
			continue
		}

		categories = append(categories, CollectionListItem{
			ColleId:   id,
			ColleName: name,
			PostCount: postCount, // 现在可以正确赋值
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
