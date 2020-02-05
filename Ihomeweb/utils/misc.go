package utils

import (
	"crypto/md5"
	"encoding/hex"
)

func AddDomain2Url(url string) (domainUrl string) {
	domainUrl = "http://" + G_fastdfs_addr + ":" + G_fastdfs_port + "/" + url
	return domainUrl
}

func GetMd5String(s string) string {
	hash := md5.New()
	hash.Write([]byte(s))
	return hex.EncodeToString(hash.Sum(nil))
}
