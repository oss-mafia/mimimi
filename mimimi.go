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
	"log"
	"regexp"
	"strings"

	"github.com/nlopes/slack"
	"github.com/nlopes/slack/slackevents"
	"github.com/pkg/errors"
)

// Bot holds the information of the Slack bot
type Bot struct {
	Username          string
	AccessToken       string
	VerificationToken string
	Api               *slack.Client
}

var vowels *regexp.Regexp

func init() {
	var err error
	vowels, err = regexp.Compile("[aeou]")
	if err != nil {
		log.Fatalf("error compiling regexp: %v\n", err)
	}
}

// New creates a new Mimimi bot
func New(username, accessToken, verificationToken string) Bot {
	return Bot{
		Username:          username,
		AccessToken:       accessToken,
		VerificationToken: verificationToken,
		Api:               slack.New(accessToken),
	}
}

// HandleMessage takes action on the given message
func (m Bot) HandleMessage(event *slackevents.MessageEvent) error {
	if event.Username != m.Username { // Do not reply to ourselves to avoid loops
		reply := strings.ToUpper(vowels.ReplaceAllString(event.Text, "I"))
		if reply != strings.ToUpper(event.Text) {
			if _, _, err := m.Api.PostMessage(event.Channel, slack.MsgOptionText(reply, false)); err != nil {
				return errors.Wrap(err, "error posting message")
			}
		}
	}
	return nil
}
