package utilities

import (
	"encoding/json"
	"errors"
	"strings"
)

// CSString is my attempt to have an array of strings in a single string and thus able to be compared for equality.
// I am not sure how much I like this as an idea, but I am going to roll with it for now. It might be cool.
// The main pitfall to this I think is that the values may themselves have commas, and if so that will cause a problem.
// I feel like I can write somthing to encode the commas inside values, but I am not 100% sold on this as a solution, so I am going to think on it for a bit.
type CSString string

var (
	ErrorCSStringIndexNegative    = errors.New("index provided must be a non negative value")
	ErrorCSStringIndexOutOfBounds = errors.New("index provided is out of bounds")
)

func NewCSString(values []string) CSString {
	var csString CSString
	for _, value := range values {
		csString.Add(value)
	}
	return csString
}

// Add adds an item to the CSString
// TODO: Stew on encoding values in the event that the contain comma themselves...
func (css *CSString) Add(value string) {
	if len(*css) == 0 {
		*css = CSString(value)
	} else {
		*css = *css + CSString(","+value)
	}
}

// Get gets an item from the CSString by its index
// will return errors if a negative index is provided or if the index is out of bounds for the string
// TODO: Stew on encoding values in the event that the contain comma themselves...
func (css *CSString) Get(index int) (string, error) {
	items := css.ToSlice()
	itemsLen := len(items)
	if index < 0 {
		return "", ErrorCSStringIndexNegative
	}
	if index > itemsLen-1 {
		return "", ErrorCSStringIndexOutOfBounds
	}
	return items[index], nil
}

func (css *CSString) Contains(value string) (int, bool) {
	for i, v := range css.ToSlice() {
		if value == v {
			return i, true
		}
	}
	return -1, false
}

func (css *CSString) ContainsCaseInsensitive(value string) (int, bool) {
	if len(*css) > 0 {
		loweredValue := strings.ToLower(value)
		for i, v := range css.ToSlice() {
			loweredV := strings.ToLower(v)
			if loweredValue == loweredV {
				return i, true
			}
		}
	}
	return -1, false
}

func (css *CSString) ToSlice() []string {
	if len(*css) == 0 {
		return nil
	}
	return strings.Split(string(*css), ",")
}

func (css *CSString) ItemCount() int {
	if len(*css) == 0 {
		return 0
	}
	return len(css.ToSlice())
}

// JSON Support

func (css CSString) MarshalJSON() ([]byte, error) {
	ss := css.ToSlice()
	data, err := json.Marshal(&ss)
	return data, err
}

func (css *CSString) UnmarshalJSON(data []byte) error {
	var ss []string
	err := json.Unmarshal(data, &ss)
	if err != nil {
		return err
	}
	*css = NewCSString(ss)
	return nil
}
