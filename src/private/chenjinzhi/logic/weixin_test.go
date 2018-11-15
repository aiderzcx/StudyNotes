package logic

import (
	"testing"
)

func TestTransWxTm2Db(t *testing.T) {

}

func TestTransDbTm2Wx(t *testing.T) {
	curTm := "2006-01-02 15:04:05"
	if "20060102150405" != transDbTm2Wx(curTm) {
		t.Errorf("not equal: %s", transDbTm2Wx(curTm))
		return
	}

}
