package client

import (
	"PAD-151-Message-Broker/model"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

//MessageParser ...
type MessageParser struct {
}

// Parse the following user commands:
// "/b [msg]" - broadcast
// "/u [username] [msg]" - send message to a specific user
// "/n [room-name]" - create a new room
// "/r [room-name] [msg]" - send message to a specific room
// "/p [msg]" - create a post
// "/s [username]" - subscribe to a user
// "/c" - get all active connections
func (mp *MessageParser) Parse(userInput string) (interface{}, error) {
	model := new(model.SentMessageModel)
	switch {
	case strings.HasPrefix(userInput, "/b "):
		re := regexp.MustCompile(`^\/b (?P<Msg>([[:graph:]]|\s)*)$`)
		msg := re.FindStringSubmatch(userInput)[1]
		model.Type = 0
		model.Message = msg

	case strings.HasPrefix(userInput, "/u "):
		re := regexp.MustCompile(`^\/u (?P<UserID>\d*) (?P<Msg>([[:graph:]]|\s)*)$`)
		strs := re.FindStringSubmatch(userInput)
		receiver, _ := strconv.Atoi(strs[1])
		model.Receivers = make([]int, 1)
		model.Receivers[0] = receiver
		model.Message = strs[2]
		model.Type = 1

	case strings.HasPrefix(userInput, "/r "):
		re := regexp.MustCompile(`^\/r (?P<RoomId>\d)*) (?P<Msg>([[:graph:]]|\s)*)$`)
		strs := re.FindStringSubmatch(userInput)
		receiver, _ := strconv.Atoi(strs[1])
		model.Receivers = make([]int, 1)
		model.Receivers[0] = receiver
		model.Message = strs[2]
		model.Type = 2

	case strings.HasPrefix(userInput, "/p "):
		re := regexp.MustCompile(`^\/p (?P<Msg>([[:graph:]]|\s)*)$`)
		msg := re.FindStringSubmatch(userInput)[1]
		model.Type = 3
		model.Message = msg

	case strings.HasPrefix(userInput, "/n "):
		re := regexp.MustCompile(`^\/n (?P<RoomName>([[:graph:]]|\s)*)$`)
		msg := re.FindStringSubmatch(userInput)[1]
		model.Type = 4
		model.Message = msg

	case strings.HasPrefix(userInput, "/s "):
		re := regexp.MustCompile(`^\/s\s+(?P<UserId>\d*)\s*$`)
		strs := re.FindStringSubmatch(userInput)
		receiver, _ := strconv.Atoi(strs[1])
		model.Type = 6
		model.Receivers = make([]int, 1)
		model.Receivers[0] = receiver

	case strings.HasPrefix(userInput, "/c"):
		if regexp.MustCompile(`^\/c\s*$`).MatchString(userInput) {
			model.Type = 7
		}

	default:
		return nil, fmt.Errorf("Its' an error")
	}
	return model, nil
}
