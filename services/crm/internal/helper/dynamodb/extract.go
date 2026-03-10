package dynamodb

import "strings"

func ExtractSK(sk string) string {
	parts := strings.SplitN(sk, "#", 2)
	if len(parts) == 2 {
		return parts[1]
	}
	return sk
}
