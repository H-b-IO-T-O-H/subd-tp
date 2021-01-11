package validation

import (
	"regexp"
	"subd/application/common/errors"
)

var patternNickname = regexp.MustCompile("^[0-9a-zA-z_.]+$")
//var pattrnSlug = regexp.MustCompile("^(\\d|\\w|-|_)*(\\w|-|_)(\\d|\\w|-|_)*$\n")

func NicknameValid(nickname string) errors.Err {
	if !patternNickname.MatchString(nickname) {
		return errors.RespErr{StatusCode: errors.BadRequestCode, Message: errors.InvalidNickname}
	}
	return nil
}
