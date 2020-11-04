package telegrambot

type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
	From From   `json:"from"`
}

type Chat struct {
	ChatID int `json:"id"`
}

type From struct {
	UserID    int    `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type RestResponse struct {
	Result []Update `json:"result"`
}

type BotMessage struct {
	ChatID      int    `json:"chat_id"`
	Text        string `json:"text"`
	ReplyMarkup string `json:"reply_markup"`
}

// type InlineKeyboard struct {
// 	InlineKeyboards []InlineKeyboard `json:"inline_keyboard"`
// }

// type InlineKeyboard struct {
// 	Text string `json:"text"`
// }
