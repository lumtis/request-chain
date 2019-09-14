package types

import "strings"

// Query Result Payload for a names query
type QueryGetBlock []string

// implement fmt.Stringer
func (n QueryGetBlock) String() string {
	return strings.Join(n[:], "\n")
}
