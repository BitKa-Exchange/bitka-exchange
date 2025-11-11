package jwt

// Claims are the custom JWT claims we use.
type Claims struct {
	UserID uint64 `json:"uid"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	// Standard claims are embedded through jwt.RegisteredClaims in implementation
}
