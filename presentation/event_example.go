type MessageWithID struct {
	Event string      `json:"event"`
	ID    int         `json:"id"`
	Data  interface{} `json:"data"` // HL
}
