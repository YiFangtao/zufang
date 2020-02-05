package models

import (
	"github.com/tedcy/fdfs_client"
	"zufang/Ihomeweb/pkg/logging"
)

//通过文件名的方式进行上传
func UploadByFilename(filename string) (RemoteFileId string, err error) {
	client, thisErr := fdfs_client.NewClientWithConfig("/Users/yift/go/zufang/Ihomeweb/conf/client.conf")
	if thisErr != nil {
		logging.Debug("上传图片失败", err)
		RemoteFileId = ""
		err = thisErr
		return
	}

	//通过句柄上传文件
	fileId, thisErr := client.UploadByFilename(filename)
	if thisErr != nil {
		logging.Debug("上传图片失败", err)
		RemoteFileId = ""
		err = thisErr
		return
	}

	//回传
	return fileId, nil
}

//功能函数 操作fast上传二进制文件
func UploadByBuffer(fileBuffer []byte, fileExtName string) (RemoteFileId string, err error) {
	//通过配置文件创建fdfs操作句柄
	client, thisErr := fdfs_client.NewClientWithConfig("/Users/yift/go/zufang/Ihomeweb/conf/client.conf")
	if thisErr != nil {
		logging.Debug("上传图片失败", thisErr)
		RemoteFileId = ""
		err = thisErr
		return
	}

	//通过句柄上传二进制文件
	fileId, thisErr := client.UploadByBuffer(fileBuffer, fileExtName)
	if thisErr != nil {
		logging.Debug("上传图片失败", err)
		RemoteFileId = ""
		err = thisErr
		return
	}

	//回传
	return fileId, nil
}
