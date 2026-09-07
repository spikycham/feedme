package model

type UserRole int

type User struct {
	ID                   int
	UserID               string
	Name                 string
	Account              string
	Password             string
	Role                 UserRole // 0 customer, 1 merchant
	AvatarURI            string
	ProfileBackgroundURI string
	CreatedAt            int64
}

const (
	UserRoleCustomer UserRole = iota
	UserRoleMerchant
)

type RefreshTokenRow struct {
	ID        int
	UserID    string
	Token     string
	CreatedAt int64
	ExpiredAt int64
}
