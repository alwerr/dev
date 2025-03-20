package dev

import "regexp"

func IsName(s string) bool {
	return len(s) > 3
}
func IsMail(s string) bool {
	match, err := regexp.MatchString(`^\w+@[a-zA-Z_]+?\.[a-zA-Z]{2,3}$`, s)
	return match && err == nil
	// e, err := mail.ParseAddress(s)
	// return err == nil && e.Address == s

}
func IsPass(s string) bool {
	return len(s) > 3
}
