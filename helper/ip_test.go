package helper_test

import (
	"github.com/Carey6918/PikaRPC/helper"
	"testing"
)

func TestGetLocalAddress(t *testing.T) {
	addr := helper.GetLocalAddress("9785")
	t.Logf("GetLocalAddress, address= %v", addr)
}
