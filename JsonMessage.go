package main

type Meta struct {
}

type User struct {
	Id       string `json:"id"`
	Nickname string `json:"nickname"`
	Login    string `json:"login"`
	Meta     Meta   `json:"meta"`
}

type JsonMessage struct {
	Id      string `json:"id"`
	Type    string `json:"type"`
	Service string `json:"service"`
	Text    string `json:"text"`
	TextSrc string `json:"text_src"`
	User    User   `json:"user"`
}
