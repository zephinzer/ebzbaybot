package types

type UserStorage map[string]User

type User struct {
	ChatID int64 `json:"chatID"`
}
