package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"runtime"
	"strings"

	"github.com/0x19/goesl"
)

// ASRText asr text
type ASRText struct {
	Text string `json:"text"`
}

// ASRResult asr 识别结果
type ASRResult struct {
	SessionID string  `json:"session_id"`
	State     string  `json:"state"`
	Version   string  `json:"version"`
	Result    ASRText `json:"result"`
}

var (
	fshost   = flag.String("fshost", "testwsssip.example.com", "Freeswitch hostname. Default: localhost")
	fsport   = flag.Uint("fsport", 8021, "Freeswitch port. Default: 8021")
	password = flag.String("pass", "ClueCon", "Freeswitch password. Default: ClueCon")
	timeout  = flag.Int("timeout", 10, "Freeswitch conneciton timeout in seconds. Default: 10")
)

// ACCBizdata 指定header
const ACCBizdata = "XCallId=MTE3Njk5NjU2NzY1XjEwNDQ5ODA5Njc2NQ=="

func doClientAPI(client *goesl.Client, jsonArgs interface{}) {
	jsonMarshal, err := json.Marshal(jsonArgs)
	if err != nil {
		goesl.Error("Error while parsing json object: %s", err)
	}

	goesl.Debug("%s", jsonMarshal)

	err = client.BgApi(fmt.Sprintf("lua lua/aiapi.lua %s", jsonMarshal))
	if err != nil {
		goesl.Error("Error while reading Freeswitch message: %s", err)
	}
}

// 呼出电话
// sessionUUID 指定呼出路的session uuid
func callOut(client *goesl.Client, sessionUUID string) {
	jsonArgs := make(map[string]interface{})
	extraHeaderArgs := make(map[string]interface{})

	extraHeaderArgs["ACC-Bizdata"] = ACCBizdata
	extraHeaderArgs["Other-Header"] = "test"

	jsonArgs["uuid"] = sessionUUID
	jsonArgs["action"] = "originate"
	jsonArgs["asr_engine"] = "puqiang"
	jsonArgs["gateway_name"] = "examplevos"
	jsonArgs["caller_number"] = "01012341234"
	jsonArgs["called_number"] = "36518513379185"
	jsonArgs["extra_headers"] = extraHeaderArgs

	doClientAPI(client, jsonArgs)
}

// 在指定session播放开场白
func playbackHello(client *goesl.Client, sessionUUID string) {
	jsonArgs := make(map[string]interface{})

	jsonArgs["uuid"] = sessionUUID
	jsonArgs["action"] = "playback"
	jsonArgs["play_type"] = "url"
	jsonArgs["arg"] = "http://pdlolj1mi.bkt.clouddn.com/b4409201810161640062358.mp3"

	doClientAPI(client, jsonArgs)
}

// 指定时间后将session挂断
// sessionUUID 待挂断session的uuid
// timeout 延时时间 单位秒
func hangupCall(client *goesl.Client, sessionUUID string, timeout int) {
	jsonArgs := make(map[string]interface{})

	jsonArgs["uuid"] = sessionUUID
	jsonArgs["time"] = timeout
	jsonArgs["action"] = "hangup"

	doClientAPI(client, jsonArgs)
}

// 播放音乐
// sessionUUID 指定session上播放声音
func playbackInfo(client *goesl.Client, sessionUUID string, url string) {
	jsonArgs := make(map[string]interface{})

	jsonArgs["uuid"] = sessionUUID
	jsonArgs["action"] = "playback"
	jsonArgs["play_type"] = "url"
	jsonArgs["arg"] = url

	doClientAPI(client, jsonArgs)
}

// 监听转接指令
// sessionUUID 指定监听路的session的uuid
// eavesdropUUID 待监听session的uuid
func transferCall(client *goesl.Client, sessionUUID string, eavesdropUUID string) {
	jsonArgs := make(map[string]interface{})

	jsonArgs["uuid"] = sessionUUID
	jsonArgs["action"] = "transfer"
	jsonArgs["eavesdrop_uuid"] = eavesdropUUID
	jsonArgs["gateway_name"] = "examplevos"
	jsonArgs["caller_number"] = "18510409326"
	jsonArgs["called_number"] = "6800000000"

	doClientAPI(client, jsonArgs)
}

func doAsrResult(client *goesl.Client, asrResult *ASRResult, transferUUID string) {
	text := asrResult.Result.Text

	if strings.Contains(text, "再见") {
		hangupCall(client, asrResult.SessionID, 0)
	} else if strings.Contains(text, "什么地方") || strings.Contains(text, "位置") {
		playbackInfo(client, asrResult.SessionID, "http://pdlolj1mi.bkt.clouddn.com/4a13f20180918184147667.wav")
	} else if strings.Contains(text, "多少钱") || strings.Contains(text, "价格") {
		playbackInfo(client, asrResult.SessionID, "http://pdlolj1mi.bkt.clouddn.com/3a2df201809181842209547.wav")
	} else if strings.Contains(text, "不需要") {
		playbackInfo(client, asrResult.SessionID, "http://pdlolj1mi.bkt.clouddn.com/2a620201809181847422822.wav")
	} else if strings.Contains(text, "转接") {
		transferCall(client, transferUUID, asrResult.SessionID)
	}
}

func main() {
	aiStop := false
	// Boost it as much as it can go ...
	runtime.GOMAXPROCS(runtime.NumCPU())

	client, err := goesl.NewClient(*fshost, *fsport, *password, *timeout)

	if err != nil {
		goesl.Error("Error while creating new client: %s", err)
		return
	}

	go client.Handle()

	err = client.Send("events json PLAYBACK_STOP CHANNEL_ANSWER CHANNEL_HANGUP_COMPLETE CUSTOM asr::answer ai::stop")
	if err != nil {
		panic(err)
	}
	// generate uuid
	//sessionUUID := uuid.NewV4().String()
	sessionUUID := "1111"
	//transferUUID := uuid.NewV4().String()
	transferUUID := "2222"

	callOut(client, sessionUUID)

RunningLoop:
	for {
		msg, err := client.ReadMessage()
		if err != nil {
			// If it contains EOF, we really dont care...
			if !strings.Contains(err.Error(), "EOF") && err.Error() != "unexpected end of JSON input" {
				goesl.Error("Error while reading Freeswitch message: %s", err)
			}
			continue
		}

		switch msg.GetHeader("Event-Name") {
		default:
		case "PLAYBACK_STOP":
			tmpSessionUUID := msg.GetHeader("variable_uuid")
			if tmpSessionUUID == sessionUUID {
				fmt.Println("music playback done")
			}
		case "CHANNEL_HANGUP_COMPLETE":
			tmpSessionUUID := msg.GetHeader("variable_uuid")
			if tmpSessionUUID == sessionUUID {
				fmt.Printf("hangup cause: %s\n", msg.GetHeader("Hangup-Cause"))
				break RunningLoop
			}
		case "CHANNEL_ANSWER":
			fmt.Println("CHANNEL_ANSWER")
			tmpSessionUUID := msg.GetHeader("variable_uuid")
			if tmpSessionUUID == sessionUUID {
				playbackHello(client, sessionUUID)
			}
		case "CUSTOM":
			goesl.Debug("%s", msg)
			switch msg.GetHeader("Event-Subclass") {
			case "ai::stop":
				fmt.Println("ai::stop")
				tmpSessionUUID := msg.GetHeader("variable_uuid")
				if tmpSessionUUID == sessionUUID {
					aiStop = true
					// break RunningLoop
				}
			case "asr::answer":
				asr := ASRResult{}
				err := json.Unmarshal(msg.Body, &asr)
				if err != nil {
					fmt.Println(err)
					break
				}
				if asr.SessionID == sessionUUID {
					fmt.Printf("customer speak text: %s\n", asr.Result.Text)
					if !aiStop {
						doAsrResult(client, &asr, transferUUID)
					}
				} else {
					fmt.Printf("other speak text: %s\n", asr.Result.Text)
				}
			default:
			}
		}
		// fmt.Printf("%s\n", msg)
	}

	fmt.Println("execute end")
}
