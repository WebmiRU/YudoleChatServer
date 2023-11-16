package main

type Config struct {
	Servers struct {
		Host struct {
			Http struct {
				Address string `json:"address"`
				Port    int16  `json:"port"`
			}
			Server struct {
				Address string `json:"address"`
				Port    int16  `json:"port"`
			}
		} `json:"host"`
	} `json:"servers"`
}

type TypeMeta struct {
	Badges map[string]string `json:"badges"` // key => Image URL
}

type TypeUser struct {
	Id       string   `json:"id"`
	Nickname string   `json:"nickname"`
	Login    string   `json:"login"`
	Meta     TypeMeta `json:"meta"`
}

type TypeMessage struct {
	Id      string   `json:"id"`
	Type    string   `json:"type"`
	Service string   `json:"service"`
	Html    string   `json:"html"`
	Text    string   `json:"text"`
	User    TypeUser `json:"user"`
}
