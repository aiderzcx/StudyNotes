/*
	安全相关的接口
*/

package wx

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"sort"
	"sync/atomic"
	"time"
)

// md5生成签名
func md5Sign(inPair map[string]string, inKey string) string {
	ks := make([]string, 0, len(inPair))

	for k := range inPair {
		if k == "sign" || k == "sign_type" || inPair[k] == "" {
			continue
		}
		ks = append(ks, k)
	}
	sort.Strings(ks)

	h := md5.New()
	signature := make([]byte, h.Size()*2)

	for i := range ks {
		h.Write([]byte(ks[i]))
		h.Write([]byte{'='})
		h.Write([]byte(inPair[ks[i]]))
		h.Write([]byte{'&'})
	}
	h.Write([]byte("key="))
	h.Write([]byte(inKey))

	hex.Encode(signature, h.Sum(nil))
	return string(bytes.ToUpper(signature))
}

// 验证md5的签名
func verifyMd5Sign(inPair map[string]string, inKey string) error {
	paramSign, exist := inPair["sign"]
	if !exist {
		return fmt.Errorf("inPair no sign: %+v", inPair)
	}

	newSign := md5Sign(inPair, inKey)
	if paramSign != newSign {
		return fmt.Errorf("paramSign(%s) != newSign(%s)", paramSign, newSign)
	}

	return nil
}

// 随机字符串，返回28位
func nonceString() string {
	return randString(28)
}

// 生产固定长度的随机字符串
var (
	g_rand_idx int64
	g_rand_str = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
)

func randString(num int) string {
	result := make([]byte, 0, num)

	randNum := time.Now().Unix()
	atomic.AddInt64(&g_rand_idx, 10)

	r := rand.New(rand.NewSource(randNum + g_rand_idx))
	for i := 0; i < num; i++ {
		result = append(result, g_rand_str[r.Intn(len(g_rand_str))])
	}

	return string(result)
}
