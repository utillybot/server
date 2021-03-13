package discord

type User struct {
	Id            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
	MfaEnabled    bool	 `json:"mfa_enabled"`
	Locale        string `json:"locale"`
	Verified      bool   `json:"verified"`
	Flags         int    `json:"flags"`
	PremiumType   int    `json:"premium_type"`
	PublicFlags   int    `json:"public_flags"`
}
