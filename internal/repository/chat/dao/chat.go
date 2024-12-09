package dao

// Chat type is the main structure for chat.
type Chat struct {
	ID        string    `db:"id"`
	Usernames []string  `db:"usernames"`
}
