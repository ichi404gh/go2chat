package main


type longpollUpdate struct {
	Ts      string         `json:"ts"`
	Updates []updateObject `json:"updates"`
}
type updateObject struct {
	Type    string  `json:"type"`
	Object  message `json:"object"`
	GroupID int     `json:"group_id"`
	Secret string `json:"secret"`
}
type failedLP struct {
	Failed int `json:"failed"`
	Ts     int `json:"ts"`
}
type message struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	FromID    int    `json:"from_id"`
	Date      int    `json:"date"`
	ReadState int    `json:"read_state"`
	Out       int    `json:"out"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	Attachments []struct {
		Type string `json:"type "`
	} `json:"attachments"`
	FwdMessages []message `json:"fwd_messages"`
}

type longpollResponse struct {
	Response *struct {
		Key    string `json:"key"`
		Server string `json:"server"`
		Ts     int    `json:"ts"`
	} `json:"response"`
	Error *struct {
		ErrorCode int    `json:"error_code"`
		ErrorMsg  string `json:"error_msg"`
		RequestParams []struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		} `json:"request_params"`
	} `json:"error"`
}