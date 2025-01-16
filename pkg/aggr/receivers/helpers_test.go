package receivers_test

import "strings"

func replaceHTTPWithWS(url string) string {
	return strings.Replace(url, "http", "ws", 1)
}
