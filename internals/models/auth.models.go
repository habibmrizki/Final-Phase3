package models

type User struct {
	ID       int     `db:"id"`
	Email    string  `db:"email"`
	Password string  `db:"password"`
	Name     string  `db:"name"`
	Avatar   *string `db:"avatar"`
	Bio      *string `db:"bio"`
}

type RegisterRequest struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required,min=6"`
	Name     string `json:"name" form:"name" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" form:"email" binding:"required,email"`
	Password string `json:"password" form:"password" binding:"required"`
}

type RegisterResponse struct {
	Message string `json:"message"`
	ID      int    `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
}

type LoginResponse struct {
	User struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"user"`
	Token string `json:"token"`
}
