package dto

type UpdateProfileRequest struct {
	FullName  string `json:"full_name"`
	AvatarURL string `json:"avatar_url"`
}
