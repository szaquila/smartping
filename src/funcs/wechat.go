package funcs

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/wzv5/pping/pkg/ping"
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
	Touser   string   `json:"touser"`
	Toparty  string   `json:"toparty"`
	Msgtype  string   `json:"msgtype"`
	Agentid  int      `json:"agentid"`
	Text     Text     `json:"text"`
	Textcard Textcard `json:"textcard"`
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

func Messages(toUser string, toParty string, agentId int, content, title, url string) string {
	msg := MESSAGES{
		Touser:  toUser,
		Toparty: toParty,
		Agentid: agentId,
		Text:    Text{Content: content},
	}
	if title != "" {
		msg.Msgtype = "textcard"
		if url == "" {
			url = "http://cacti197.a.yjidc.com:8899"
		}
		msg.Textcard = Textcard{Title: title, Description: content, Url: url}
	} else {
		msg.Msgtype = "text"
		msg.Text = Text{Content: content}
	}
	sed_msg, _ := json.Marshal(msg)
	//fmt.Printf("%s",string(sed_msg))
	return string(sed_msg)
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

func StringToInt(e string) (int, error) {
	return strconv.Atoi(e)
}

func StringToInt64(e string) (int64, error) {
	return strconv.ParseInt(strings.Trim(e, " "), 10, 64)
}

func Int64ToString(e int64) string {
	return strconv.FormatInt(e, 10)
}

func PingToChan(ctx context.Context, p ping.IPing) <-chan ping.IPingResult {
	c := make(chan ping.IPingResult)
	go func() {
		c <- p.PingContext(ctx)
	}()
	return c
}

func RunPing(p ping.IPing) (Statistics, error) {
	count := 5
	interval := time.Second * 1
	if count > 1 {
		// 预热，由于某些资源需要初始化，首次运行会耗时较长
		p.Ping()
	}

	s := Statistics{}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())

	for i := 1; i <= count; i++ {
		select {
		case result := <-PingToChan(ctx, p):
			PrintResult(i, result)
			s.append(result)
		case <-c:
			goto end
		}

		// 最后一次 ping 结束后不再等待
		if i == count {
			break
		}

		select {
		case <-c:
			goto end
		case <-time.After(interval):
		}

		// 再次检查是否停止，上面的检查可能由于延迟为0而始终无法停止
		select {
		case <-c:
			goto end
		default:
		}
	}

end:
	cancel()
	// if count > 1 {
	// 	s.print()
	// }
	if s.Sent == 0 || s.Failed != 0 {
		return s, errors.New("ping error")
	}
	return s, nil
}

type Statistics struct {
	Max, Min, Total  float64
	Sent, Ok, Failed int
}

func (s *Statistics) append(result ping.IPingResult) {
	if result == nil {
		return
	}
	s.Sent++
	if result.Error() != nil {
		s.Failed++
		return
	}
	t := float64(result.Result())
	if s.Ok == 0 {
		s.Min = t
		s.Max = t
	} else {
		if t < s.Min {
			s.Min = t
		} else if t > s.Max {
			s.Max = t
		}
	}
	s.Total += t
	s.Ok++
}

func (s *Statistics) clear() {
	s.Max = 0
	s.Min = 0
	s.Total = 0
	s.Sent = 0
	s.Ok = 0
	s.Failed = 0
}

func (s *Statistics) print() {
	if s.Sent == 0 {
		return
	}
	fmt.Println()
	fmt.Printf("\tsent = %d, ok = %d, failed = %d (%.2f%%)\n", s.Sent, s.Ok, s.Failed, float64(100*s.Failed)/float64(s.Sent))
	if s.Ok > 0 {
		fmt.Printf("\tmin = %f ms, max = %f ms, avg = %.2f ms\n", s.Min, s.Max, s.Total/float64(s.Ok))
	}
}

func PrintResult(i int, r ping.IPingResult) {
	log.Printf("[%d] %v\n", i, r)
}
