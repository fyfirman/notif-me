package telegram

type Message struct {
	MessageID int `json:"message_id"`
	From      struct {
		ID           int    `json:"id"`
		IsBot        bool   `json:"is_bot"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		Username     string `json:"username"`
		LanguageCode string `json:"language_code"`
	} `json:"from"`
	Chat struct {
		ID                          int64  `json:"id"`
		Title                       string `json:"title"`
		Type                        string `json:"type"`
		AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
	} `json:"chat"`
	Date     int64  `json:"date"`
	Text     string `json:"text"`
	Entities []struct {
		Offset int    `json:"offset"`
		Length int    `json:"length"`
		Type   string `json:"type"`
	} `json:"entities"`
}

type OnUpdateMessageBody struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}
