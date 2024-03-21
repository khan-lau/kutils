package ktest

import (
	"fmt"
	"testing"

	"github.com/khan-lau/kutils/kuuid"
)

func TestUuidV1(t *testing.T) {
	v1, err := kuuid.NewV1()
	u1 := kuuid.Must(v1, err)
	str := u1.String()
	fmt.Println()
	fmt.Printf("UUIDv1: %s\n", u1)
	timestamp, err := kuuid.GetTimestampFromUUIDV1String(str)
	if err != nil {
		t.Errorf("%s", err)
		fmt.Printf("GetTimestampFromUuidV1 error: %s\n", err)
	}
	fmt.Printf("timestamp: %s\n", timestamp)

	timestamp = kuuid.GetTimestampFromUUIDV1(u1)
	fmt.Printf("timestamp: %s\n", timestamp)

	u1, err = kuuid.UUIDv1FromString(str)
	if err != nil {
		t.Errorf("%s", err)
		fmt.Printf("UUIDv1FromString error: %s\n", err)
	}
	timestamp = kuuid.GetTimestampFromUUIDV1(u1)
	fmt.Printf("timestamp: %s\n", timestamp)
}
