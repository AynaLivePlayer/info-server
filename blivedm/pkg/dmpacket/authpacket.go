package dmpacket

import "encoding/json"

func GenerateAuthBody(uid int, roomId int, token string) string {
	data := map[string]interface{}{
		"uid":      uid,
		"roomid":   roomId,
		"protover": 2,
		"platform": "web",
		"type":     2,
		"key":      token,
	}
	val, _ := json.Marshal(data)
	return string(val)
}
