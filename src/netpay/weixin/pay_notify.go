package wx

import (
	"encoding/xml"
)

type PayNotifyResp struct {
	XMLName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code"`
	ReturnMsg  string   `xml:"return_msg"`
}

func PayNotify(inPayCfg ConfItemST, inParam map[string]string, outResp *PayNotifyResp) error {
	return nil
}
