package auth

type GoogleAPIUserinfoRes struct {
	ID            string `json:"id" binding:"required"`
	Email         string `json:"email" binding:"required"`
	VerifiedEmail bool   `json:"verified_email" binding:"required"`
	Name          string `json:"name" binding:"required"`
	GivenName     string `json:"given_name" binding:"required"`
	FamilyName    string `json:"family_name" binding:"required"`
	Picture       string `json:"picture" binding:"required"`
	Locale        string `json:"locale" binding:"required"`
}
