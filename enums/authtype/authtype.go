package authtype

type UserAuthType uint

const (
	Local       UserAuthType = 1
	GoogleOAuth UserAuthType = 2
)
