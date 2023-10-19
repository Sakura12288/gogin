package models

//验证token对应的是否存在

type Auth struct {
	ID       int    `gorm:"primary_key" json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func CheckAuth(username, password string) (bool, error) {
	var auth Auth
	if err := db.Select("id").Where(Auth{Username: username, Password: password}).First(&auth).Error; err != nil {
		return false, err
	}
	if auth.ID > 0 {
		return true, nil
	}
	return false, nil
}
