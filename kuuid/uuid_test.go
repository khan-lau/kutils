package kuuid

import (
	"fmt"
	"testing"
)

func TestUuidV1(t *testing.T) {
	v1, err := NewV1()
	u1 := Must(v1, err)
	str := u1.String()
	fmt.Println()
	fmt.Printf("UUIDv1: %s\n", u1)
	timestamp, err := GetTimestampFromUUIDV1String(str)
	if err != nil {
		t.Errorf("%s", err)
		fmt.Printf("GetTimestampFromUuidV1 error: %s\n", err)
	}
	fmt.Printf("timestamp: %s\n", timestamp)

	timestamp = GetTimestampFromUUIDV1(u1)
	fmt.Printf("timestamp: %s\n", timestamp)

	u1, err = UUIDv1FromString(str)
	if err != nil {
		t.Errorf("%s", err)
		fmt.Printf("UUIDv1FromString error: %s\n", err)
	}
	timestamp = GetTimestampFromUUIDV1(u1)
	fmt.Printf("timestamp: %s\n", timestamp)
}
