package main

import "encoding/json"

// Request .
type Request struct {
	Meta struct {
		Locale     string `json:"locale"`
		Timezone   string `json:"timezone"`
		ClientID   string `json:"client_id"`
		Interfaces struct {
			Screen struct{} `json:"screen"`
		} `json:"interfaces"`
	} `json:"meta"`

	Request struct {
		Command string `json:"command"`
		NLU     struct {
			Tokens   []string `json:"tokens"`
			Entities []struct {
				Tokens struct {
					Start int `json:"start"`
					End   int `json:"end"`
				} `json:"tokens"`
				Type string `json:"type"`
				// Value NEType `json:"value"`
				Value *json.RawMessage `json:"value"`
			} `json:"entities"`
		} `json:"nlu"`
		Utterance string `json:"original_utterance"`
		Type      string `json:"type"`
		Markup    struct {
			DangerousContext bool `json:"dangerous_context,omitempty"`
		} `json:"markup,omitempty"`
		Payload interface{} `json:"payload,omitempty"`
	} `json:"request"`

	Session struct {
		New       bool   `json:"new"`
		MessageID int    `json:"message_id"`
		SessionID string `json:"session_id"`
		SkillID   string `json:"skill_id"`
		UserID    string `json:"user_id"`
	} `json:"session"`

	Version string `json:"version"`
}

// Response .
type Response struct {
	Response struct {
		Text string            `json:"text"`
		TTS  string            `json:"tts,omitempty"`
		Card ResponseWithImage `json:"card,omitempty"`
		// Buttons []struct {
		// 	Title   string   `json:"title"`
		// 	Payload interface{} `json:"payload,omitempty"`
		// 	URL     string   `json:"url,omitempty"`
		// 	Hide    bool     `json:"hide"`
		// } `json:"buttons,omitempty"`
		Buttons    []Button `json:"buttons,omitempty"`
		EndSession bool     `json:"end_session"`
	} `json:"response"`

	Session struct {
		MessageID int    `json:"message_id"`
		SessionID string `json:"session_id"`
		UserID    string `json:"user_id"`
	} `json:"session"`

	Version string `json:"version"`
}

const (
	// BigImageType .
	BigImageType = "BigImage"
	// ItemsListType .
	ItemsListType = "ItemsList"
)

// BigImage .
type BigImage struct {
	Type        string `json:"type"`
	ImageID     string `json:"image_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Button      Button `json:"button,omitempty"`
}

// ItemsList .
type ItemsList struct {
	Type   string `json:"type"`
	Header struct {
		Text string `json:"text"`
	} `json:"header"`
	Items []struct {
		ImageID     string `json:"image_id"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Button      Button `json:"button,omitempty"`
	} `json:"items"`
	Footer struct {
		Text   string `json:"text"`
		Button struct {
			Text    string      `json:"text"`
			URL     string      `json:"url,omitempty"`
			Payload interface{} `json:"payload,omitempty"`
		} `json:"button"`
	} `json:"footer"`
}

func (*BigImage) image()  {}
func (*ItemsList) image() {}

// ResponseWithImage .
type ResponseWithImage interface {
	image()
}

// Button .
type Button struct {
	Title   string      `json:"title"`
	Payload interface{} `json:"payload,omitempty"`
	URL     string      `json:"url,omitempty"`
	Hide    bool        `json:"hide,omitempty"`
}

const (
	// CommandTypeVoice .
	CommandTypeVoice = "SimpleUtterance"
	// CommandTypeButton .
	CommandTypeButton = "ButtonPressed"
)
