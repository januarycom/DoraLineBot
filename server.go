// Copyright 2016 LINE Corporation
//
// LINE Corporation licenses this file to you under the Apache License,
// version 2.0 (the "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at:
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package main

import (
	"log"
	"net/http"
	"time"
	"fmt"
	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	loc, _ := time.LoadLocation("Asia/Tokyo")
	hh := []int{6,6,6,6,6,6,7,7,7,7,7,7,7,7,7,7,7,7,7,7,7,7,7,7,7,7,7,8,8,8,8,8,8,8,8,8,8,8,8,8,9,9,9,9,9,9,9,9,9,9,9,9,9,9,10,10,10,10,10,10,10,10,10,10,10,10,10,10,10,11,11,11,11,11,11,11,11,11,11,11,11,11,11,11,12,12,12,12,12,12,12,12,12,12,12,12,12,12,12,13,13,13,13,13,13,13,13,13,13,13,13,13,14,14,14,14,14,14,14,14,14,14,14,14,14,15,15,15,15,15,15,15,15,15,15,15,15,15,16,16,16,16,16,16,16,16,16,16,16,16,16,16,16,17,17,17,17,17,17,17,17,17,17,17,17,17,17,17,17,17,17,17,17,17,18,18,18,18,18,18,18,18,18,18,18,18,18,18,18,18,18,19,19,19,19,19,19,19,19,19,20,20,20,20,20,20,21,21,21,21,21,21}
	mm := []int{29,39,44,49,56,59,4,6,14,17,19,22,24,25,27,29,32,34,37,39,42,44,47,49,52,54,57,4,7,14,19,24,27,32,34,37,42,47,52,57,2,7,12,16,19,22,29,31,34,37,40,43,49,55,0,4,7,12,16,19,22,29,31,34,37,40,43,49,55,0,3,5,7,12,19,22,29,31,34,37,40,43,49,55,1,4,7,12,15,19,22,29,31,34,37,40,43,49,55,1,7,12,19,22,25,29,31,34,40,43,49,55,1,4,7,12,19,22,29,31,34,37,43,46,55,1,4,7,9,16,19,24,31,34,41,43,51,55,1,7,9,16,19,24,29,35,39,42,44,49,51,54,59,3,4,9,14,17,19,23,29,30,30,32,34,36,39,42,44,47,49,52,54,59,1,4,7,9,14,17,19,24,27,29,30,32,34,39,42,48,52,1,5,13,15,20,25,28,41,53,3,15,24,35,44,58,4,15,26,40,48,56}

	bot, err := linebot.New(
		"ChannelSecret",
		"AccessToken")
	if err != nil {
		log.Fatal(err)
	}

	// Setup HTTP Server for receiving requests from LINE platform
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					//メッセージ生成処理
					nt := time.Now()

					reply := ""
					for i, v := range hh {
						d := time.Date(nt.Year(), nt.Month(), nt.Day(), v, mm[i], 0, 0, loc)

						if nt.Before(d){
						//if nt.Before(nt.Add(d.Sub(nt))){
							sec := d.Sub(nt).Seconds()

							if sec < 60{
								//1分切っているとき
								reply = fmt.Sprintf("%d秒で支度しな!!\n(次のバスの到着時刻　%d:%d)", int(sec), hh[i], mm[i])
							}else if sec < 3600 {
								//1時間切っているとき
								reply = fmt.Sprintf("%d分で支度しな!!\n(次のバスの到着時刻　%d:%d)", int(sec/60), hh[i], mm[i])
							} else {
								//1時間以上
								reply = fmt.Sprintf("%d時間%d分で支度しな!!\n(次のバスの到着時刻　%d:%d)", int(sec/3600), (int(sec)%3600)/60, hh[i], mm[i])
							}

							break
						}
					}

					if len(reply) < 1 {
						//LineBot応答処理
						imgurl := "https://pbs.twimg.com/media/Cj2mP8sWYAAEHov.jpg:small"
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewImageMessage(imgurl, imgurl)).Do(); err != nil {
							log.Print(err)
						}
					} else {
						//LineBot応答処理
						if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(reply)).Do(); err != nil {
							log.Print(err)
						}
					}

					//メッセージ生成ここまで

					log.Print(message)
				}
			}

		}
	})
	// This is just sample code.
	// For actual use, you must support HTTPS by using `ListenAndServeTLS`, a reverse proxy or something else.
	if err := http.ListenAndServe(":1337", nil); err != nil {
		log.Fatal(err)
	}
}
