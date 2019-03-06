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

package mimimi

import (
	"testing"

	"github.com/nlopes/slack/slackevents"
)

func TestMimimize(t *testing.T) {
	values := map[string]string{
		"":                "",
		"Hi, there!":      "HI, THIRI!",
		"mimimi123":       "MIMIMI123",
		"HELLO oss-mafia": "HILLI ISS-MIFII",
	}

	for k, v := range values {
		if res := mimimize(k); res != v {
			t.Fatalf("mimimi(%s) = %q, expected %q", k, res, v)
		}
	}
}

func TestDoNotSelfReply(t *testing.T) {
	replyCount := 0
	bot := Bot{
		Username: "mimimi",
		reply: func(target string, message string) error {
			replyCount++
			return nil
		},
	}

	if err := bot.HandleMessage(&slackevents.MessageEvent{Username: "mimimi"}); err != nil {
		t.Fatalf("unexpected error handling message: %v", err)
	}

	if replyCount > 0 {
		t.Fatalf("expected bot to not reply to itself")
	}
}

func TestDoNotReplySameMessage(t *testing.T) {
	replyCount := 0
	bot := Bot{
		Username: "mimimi",
		reply: func(target string, message string) error {
			replyCount++
			return nil
		},
	}

	if err := bot.HandleMessage(&slackevents.MessageEvent{Username: "foo", Text: "MIMIMI"}); err != nil {
		t.Fatalf("unexpected error handling message: %v", err)
	}

	if replyCount > 0 {
		t.Fatalf("expected bot to not reply if the message is already mimimized")
	}
}

func TestHandleMessage(t *testing.T) {
	reply := ""
	bot := Bot{
		Username: "mimimi",
		reply: func(target string, message string) error {
			reply = message
			return nil
		},
	}

	message := "Hi there!"
	if err := bot.HandleMessage(&slackevents.MessageEvent{Username: "foo", Text: message}); err != nil {
		t.Fatalf("unexpected error handling message: %v", err)
	}

	if reply != mimimize(message) {
		t.Fatalf("expected bot to have replied with the mimimized message")
	}
}
