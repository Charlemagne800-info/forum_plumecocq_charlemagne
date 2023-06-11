package main

import (
	"Forum_Ynov/backend/handlers/VerfiUser"
	"testing"
)

func main() {
	VerfiUser.VerifUser()
	TestVerifUser()

}

func TestVerifUser(t *testing.T) {
	VerifUser := VerifUser("test", "test")
	if VerifUser == true {
		t.Errorf("VerifUser() = %v; want false", VerifUser)
	}
}
