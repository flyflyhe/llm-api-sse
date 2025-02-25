package form

type DsRequest struct {
	ChatSessionId   string  `json:"chat_session_id"`
	ParentMessageId int32   `json:"parent_message_id"`
	Prompt          string  `json:"prompt"`
	RefFileIds      []int32 `json:"ref_file_ids"`
	ThinkingEnabled bool    `json:"thinking_enabled"`
	SearchEnabled   bool    `json:"search_enabled"`
}
