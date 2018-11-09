package wx

import (
	"testing"
)

func init() {
	g_conf.Scan.AppId = "scan_appid"
	g_conf.Scan.MchId = "scan_mchid"
	g_conf.Scan.SecretKey = "scanSecrityKey"
}

func TestNonceString(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Logf(nonceString())
	}
}

func TestMd5Sign(t *testing.T) {
	params := map[string]string{
		"appid":  "11111",
		"mch_id": "22222",
		"aaaa":   "aaaaa",
		"bbbb":   "bbbbbb",
		"cccc":   "cccccc",
	}

	secKey := "abcdefghijklmn"
	rightSign := "3EB2503CDCA56E85EC49DDCE10A2C30F"

	newSign := md5Sign(params, secKey)
	if newSign == rightSign {
		t.Logf("newSign(%s) is right", newSign)
	} else {
		t.Errorf("newSign(%s) is error", newSign)
	}
}

func TestCreateCodeLink(t *testing.T) {
	link := CreateCodeLink("product_id")
	t.Logf("CreateCodeLink: %s", link)
}
