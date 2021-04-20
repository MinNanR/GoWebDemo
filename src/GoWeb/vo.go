package main

type LoginVO struct {
	RedirectUrl string   `json:"redirectUrl"`
	JwtToken    string   `json:"jwtToken"`
	Router      []string `json:"router"`
}

//type RouterList []string
//
//func (rl RouterList) MarshalJSON() ([]byte, error) {
//	list := strings.Join(rl, ",")
//	return []byte(list), nil
//}
