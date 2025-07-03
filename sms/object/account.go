package object

type Account struct {
	ID       int    `xorm:"'ID' pk autoincr"`
	Fullname string `xorm:"'fullname'"`
	Email    string `xorm:"'email'"`
	Username string `xorm:"'username'"`
	Password string `xorm:"'password'"`
}
