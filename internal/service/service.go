package service

import "database/sql"

type Services struct {
	Beans    *BeanService
	Roasters *RoasterService
	Users    *UserService // interacts with permissions
}

func NewServices(db *sql.DB) *Services {
	return &Services{
		Beans:    NewBeanService(db),
		Roasters: NewRoasterService(db),
		Users:    NewUserService(db),
	}
}
