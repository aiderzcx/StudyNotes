/*
相关的接口参数定义
*/

package wapi

import (
	"fmt"
)

type ScanLinkReq struct {
	ProductId string `form:"product_id"` // 产品ID
}

func (s *ScanLinkReq) ParamCheck() error {
	if s.ProductId == "" {
		return fmt.Errorf("product_id is nil")
	}

	return nil
}

type ScanLinkResp struct {
	Link string `json:"link"`
}
