package account

import "time"

type CreateAccount struct {
	ParentID string `json:"parentId,omitempty"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Basis    string `json:"basis"`
}

type Account struct {
	ID        string    `json:"id"`
	ParentID  string    `json:"parentId,omitempty" `
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Basis     string    `json:"basis"`
	CreatedAt time.Time `json:"createdAt"`
}
