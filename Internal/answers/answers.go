package answers

import (
	"errors"
	"net/http"
)

var AnPreviousMessageNotRead = errors.New("the previous message was not read")
var AnEmpty = errors.New("empty, nothing found")
var AnPudgeNotLose = errors.New("the pudge not lose weight")

func StatusHTTP(e error) int {
	switch e {
	case AnPreviousMessageNotRead:
		return http.StatusConflict
	case AnEmpty:
		return http.StatusNoContent
	default:
		return http.StatusOK
	}
}
