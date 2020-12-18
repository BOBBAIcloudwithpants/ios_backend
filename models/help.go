package models

func CreateHelp(forum_id int, user_id int, title string, content string, bonus int, filename string) (int64, error) {
	sql :=
		`
			INSERT INTO help(forum_id, user_id, title, content, bonus, filename) VALUES(?,?,?,?,?,?);
		`
	return Execute(sql, forum_id, user_id, title, content, bonus, filename)
}
