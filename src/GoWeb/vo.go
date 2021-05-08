package main

type LoginVO struct {
	RedirectUrl string   `json:"redirectUrl"`
	JwtToken    string   `json:"jwtToken"`
	Router      []string `json:"router"`
}

type DownloadToolsVO struct {
	DownloadUrl string `json:"downloadUrl"`
	FileName    string `json:"fileName"`
}

//type RouterList []string
//
//func (rl RouterList) MarshalJSON() ([]byte, error) {
//	list := strings.Join(rl, ",")
//	return []byte(list), nil
//}
