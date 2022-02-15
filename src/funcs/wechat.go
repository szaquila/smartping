package funcs

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type TOKEN struct {
	ErrCode     int64     `json:"errcode"`
	ErrMsg      string    `json:"errmsg"`
	AccessToken string    `json:"access_token"`
	ExpiresIn   int64     `json:"expires_in"`
	Time        time.Time `json:"-"`
}

type Text struct {
	Content string `json:"content"`
}

type Textcard struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Btntxt      string `json:"btntxt,omitempty"`
}

type MESSAGES struct {
	Touser  string `json:"touser"`
	Toparty string `json:"toparty"`
	Msgtype string `json:"msgtype"`
	Agentid int    `json:"agentid"`
	Text    Text   `json:"text"`
	Safe    int    `json:"safe"`
}

func GetAccessToken(corpId, corpSecret string) TOKEN {
	gettoken_url := "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=" + corpId + "&corpsecret=" + corpSecret
	//print(gettoken_url)
	client := &http.Client{}
	req, err := client.Get(gettoken_url)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	body, _ := ioutil.ReadAll(req.Body)
	//fmt.Printf("\n%q",string(body))
	var json_str TOKEN
	json.Unmarshal([]byte(body), &json_str)
	json_str.Time = time.Now()
	//fmt.Printf("\n%q",json_str.Access_token)
	return json_str
}

func SendMessage(accessToken, msg string) (err error) {
	send_url := "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=" + accessToken
	//print(send_url)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", send_url, bytes.NewBuffer([]byte(msg)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("charset", "UTF-8")
	resp, err := client.Do(req)
	if err == nil {
		defer resp.Body.Close()
	}
	// panic(err)
	//fmt.Println("response Status:", resp.Status)
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println("response Body:", string(body))
	return
}

func Messages(toUser string, toParty string, agentId int, content string) string {
	msg := MESSAGES{
		Touser:  toUser,
		Toparty: toParty,
		Msgtype: "textcard",
		Agentid: agentId,
		Safe:    0,
		Text:    Text{Content: content},
	}
	sed_msg, _ := json.Marshal(msg)
	//fmt.Printf("%s",string(sed_msg))
	return string(sed_msg)
}

func StringToInt(e string) (int, error) {
	return strconv.Atoi(e)
}

func StringToInt64(e string) (int64, error) {
	return strconv.ParseInt(strings.Trim(e, " "), 10, 64)
}

func Int64ToString(e int64) string {
	return strconv.FormatInt(e, 10)
}

// func main() {
// 	touser := "BigBoss"  //企业号中的用户帐号，在zabbix用户Media中配置，如果配置不正常，将按部门发送。
// 	toparty := "2"       //企业号中的部门id。
// 	agentid:= 1000002    //企业号中的应用id。
// 	corpid := "xxxxxxxxxxxxxxxxx"  //企业号的标识
// 	corpsecret := "exxxxxxxxxxxxxxxxxxxxxxxxxxxxx"  ///企业号中的应用的Secret
// 	accessToken := Get_AccessToken(corpid,corpsecret)
// 	content := os.Args[1]
// 	//  fmt.Println(content)
// 	// 序列化成json之后，\n会被转义，也就是变成了\\n，使用str替换，替换掉转义
// 	msg := strings.Replace(messages(touser,toparty,agentid,content),"\\\\","\\",-1)

// 	//  fmt.Println(strings.Replace(msg,"\\\\","\\",-1))
// 	Send_Message(accessToken,msg)
// }
