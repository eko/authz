package configs

type User struct {
	AdminDefaultPassword string `config:"user_admin_default_password"`
}

func newUser() *User {
	return &User{
		AdminDefaultPassword: "changeme",
	}
}
