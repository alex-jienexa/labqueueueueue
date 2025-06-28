package models

// Student является структурой студента-пользователя системы.
// Admin = староста, только старосты могут создавать очередь и менять
// её параметры.
type Student struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Username string `json:"username"`
	Password string `json:"-"` // Пароль не возвращаем в JSON
	IsAdmin  bool   `json:"is_admin"`
}
