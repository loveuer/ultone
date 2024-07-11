package sqlType

import "encoding/json"

type Set map[string]struct{}

func (s Set) MarshalJSON() ([]byte, error) {
	array := make([]string, 0)
	for name := range s {
		array = append(array, name)
	}

	return json.Marshal(array)
}

func (s *Set) UnmarshalJSON(b []byte) error {
	array := make([]string, 0)
	if err := json.Unmarshal(b, &array); err != nil {
		return err
	}

	set := make(map[string]struct{})

	for _, name := range array {
		set[name] = struct{}{}
	}

	*s = set
	return nil
}

func (s Set) ToStringSlice() []string {
	var (
		result = make([]string, 0, len(s))
	)

	for key := range s {
		result = append(result, key)
	}

	return result
}

func (s *Set) FromStringSlice(ss *[]string) {
	if s == nil {
		m := make(Set)
		s = &m
	}

	for idx := range *(ss) {
		(*s)[(*ss)[idx]] = struct{}{}
	}
}
