package ostrovok

// OstrovokUser представляет пользователя из системы Островок
type OstrovokUser struct {
	Login string `json:"login"`
	Email string `json:"email"`
}
