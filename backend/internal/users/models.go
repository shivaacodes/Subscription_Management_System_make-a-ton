package users

import (
	"time"
)

type Role string

const (
	RoleAdmin    Role = "ADMIN"
	RoleStaff    Role = "STAFF"
	RoleCustomer Role = "CUSTOMER"
)

func (r Role) IsValid() bool {
	switch r {
	case RoleAdmin, RoleStaff, RoleCustomer:
		return true
	default:
		return false
	}
}

type User struct {
	ID           string    `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type CustomerProfile struct {
	ID             string    `json:"id"`
	UserID         string    `json:"userId"`
	CompanyName    *string   `json:"companyName,omitempty"`
	Phone          *string   `json:"phone,omitempty"`
	BillingAddress *string   `json:"billingAddress,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
