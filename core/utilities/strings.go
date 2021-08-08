package utilities

import "strings"

func UpperCaseFirstChar(input string) string {
	length := len(input)
	if length == 0 {
		return input
	}
	inputBuffer := make([]byte, length)
	for i, c := range input {
		var char byte = byte(c)
		if i == 0 {
			char = strings.ToUpper(string(char))[0]
		}
		inputBuffer[i] = string(char)[0]
	}
	return string(inputBuffer)
}

func LowerCaseFirstChar(input string) string {
	length := len(input)
	if length == 0 {
		return input
	}
	inputBuffer := make([]byte, length)
	for i, c := range input {
		var char byte = byte(c)
		if i == 0 {
			char = strings.ToLower(string(char))[0]
		}
		inputBuffer[i] = string(char)[0]
	}
	return string(inputBuffer)
}
