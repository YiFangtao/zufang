package utils

func AddDomain2Url(url string) (domainUrl string) {
	domainUrl = "http://" + G_fastdfs_addr + ":" + G_fastdfs_port + "/" + url
	return domainUrl
}
