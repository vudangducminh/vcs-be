package entities

type Account struct {
	ID       int    `xorm:"'id' pk autoincr"`
	Fullname string `xorm:"'fullname'"`
	Email    string `xorm:"'email'"`
	Username string `xorm:"'username'"`
	Password string `xorm:"'password'"`
	Role     string `xorm:"'role'"`
}
