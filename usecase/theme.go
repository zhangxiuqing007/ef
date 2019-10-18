package usecase

import "ef/models"

//QueryAllThemes 获取所有的主题指针
func QueryAllThemes() ([]*models.ThemeInDB, error) {
	//先从数据库读取
	return db.QueryAllThemes()
}

//QueryTheme 获取主题
func QueryTheme(themeID int) (*models.ThemeInDB, error) {
	return db.QueryTheme(themeID)
}

//QueryPostCountOfTheme 查询主题的帖子量
func QueryPostCountOfTheme(themeID int) (int, error) {
	return db.QueryPostCountOfTheme(themeID)
}
