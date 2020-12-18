package models

import "strconv"

// 在数据库中存放的 File 类型的定义
type ExtendedFile struct {
	FileID   int    `json:"file_id"`
	PostID   int    `json:post_id`
	FileName string `json:"filename"`
	Bucket   string `json:"bucket"`
	CreateAt string `json:"create_at"`
}

// 将数据库查询结果转换为 Forum
func convertMapToExtendedFile(extendedFile map[string]string) ExtendedFile {
	file_id, _ := strconv.Atoi(extendedFile["file_id"])
	post_id, _ := strconv.Atoi(extendedFile["post_id"])
	return ExtendedFile{
		FileID:   file_id,
		PostID:   post_id,
		FileName: extendedFile["filename"],
		Bucket:   extendedFile["bucket"],
		CreateAt: extendedFile["create_at"],
	}
}

// 创建文件记录
func CreateFile(file ExtendedFile) (int64, error) {
	sentence := "INSERT INTO file(post_id, filename, bucket) VALUES(?, ?, ?)"
	return Execute(sentence, file.PostID, file.FileName, file.Bucket)
}

// 根据 post_id 获取其相关的 files
func GetFilesByPostID(post_id int) ([]ExtendedFile, error) {
	var files []ExtendedFile

	sentence :=
		`SELECT * FROM file WHERE post_id = ?`
	res, err := QueryRows(sentence, post_id)

	if err != nil {
		return nil, err
	}

	for _, r := range res {
		files = append(files, convertMapToExtendedFile(r))
	}

	return files, nil
}
