// Copyright 2019 The OSS Mafia team
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/nlopes/slack/slackevents"
	"github.com/pkg/errors"

	"github.com/oss-mafia/mimimi"
)

func main() {
	accessToken := os.Getenv("ACCESS_TOKEN")
	verificationToken := os.Getenv("VERIFICATION_TOKEN")
	if accessToken == "" || verificationToken == "" {
		log.Fatal("the ACCESS_TOKEN and VERIFICATION_TOKEN env vars are required")
	}
	bot := mimimi.New("mimimi", accessToken, verificationToken)

	http.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
		event, err := verify(bot.VerificationToken, w, r)
		if err != nil {
			log.Printf("[error] %v\n", err)
			return
		}

		if event.Type == slackevents.CallbackEvent {
			innerEvent := event.InnerEvent
			switch ev := innerEvent.Data.(type) {
			case *slackevents.MessageEvent:
				if err := bot.HandleMessage(ev); err != nil {
					log.Printf("[error] %v\n", err)
				}
			}
		}
	})

	var port int
	flag.IntVar(&port, "port", 4390, "port where the bot will listen for events")
	flag.Parse()

	log.Printf("[info] listening on port %d", port)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal(err)
	}
}

// verify the incoming request and return the parsed event.  If the request is the specific URL challenge validation,
// this method replies with he challenge so callers just have to deal with regular EventAPI events
func verify(verificationToken string, w http.ResponseWriter, r *http.Request) (slackevents.EventsAPIEvent, error) {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(r.Body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return slackevents.EventsAPIEvent{}, errors.Wrap(err, "error reading message body")
	}
	body := buf.String()

	options := slackevents.OptionVerifyToken(&slackevents.TokenComparator{VerificationToken: verificationToken})
	event, err := slackevents.ParseEvent(json.RawMessage(body), options)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return slackevents.EventsAPIEvent{}, errors.Wrap(err, "error parsing event")
	}

	if event.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return slackevents.EventsAPIEvent{}, errors.Wrap(err, "error parsing challenge")
		}

		w.Header().Set("Content-Type", "text")
		if _, err = w.Write([]byte(r.Challenge)); err != nil {
			return slackevents.EventsAPIEvent{}, errors.Wrap(err, "error replying with challenge")
		}
	}

	return event, nil
}
