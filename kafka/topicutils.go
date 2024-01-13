package kafka

import "strings"

func TopicFromParts(parts ...string) string {
	return strings.Join(parts, ".")
}
