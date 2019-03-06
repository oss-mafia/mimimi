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
	API               *slack.Client

	// reply function to be used to post messages back. Useful to allow changing it for
	// testing
	reply func(target string, message string) error
}

var vowels *regexp.Regexp

func init() {
	var err error
	vowels, err = regexp.Compile("[aeouAEOU]")
	if err != nil {
		log.Fatalf("error compiling regexp: %v\n", err)
	}
}

// New creates a new Mimimi bot
func New(username, accessToken, verificationToken string) Bot {
	api := slack.New(accessToken)
	return Bot{
		Username:          username,
		AccessToken:       accessToken,
		VerificationToken: verificationToken,
		API:               api,
		reply: func(target string, message string) error {
			_, _, err := api.PostMessage(target, slack.MsgOptionText(message, false))
			return err
		},
	}
}

// HandleMessage takes action on the given message
func (m Bot) HandleMessage(event *slackevents.MessageEvent) error {
	if event.Username != m.Username { // Do not reply to ourselves to avoid loops
		reply := mimimize(event.Text)
		if reply != strings.ToUpper(event.Text) {
			if err := m.reply(event.Channel, reply); err != nil {
				return errors.Wrap(err, "error posting message")
			}
		}
	}
	return nil
}

// mimimize the given text
func mimimize(text string) string {
	return strings.ToUpper(vowels.ReplaceAllString(text, "I"))
}
