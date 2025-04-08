package utils

import (
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/objects"
)

// BatchUpdateFileLifecycle 批量更新图片的有效期为永久
func BatchUpdateFileLifecycle(fileNames []string, cover, accessKey, secretKey, bucketName string) error {
	// 构建 Bucket 对象
	mac := credentials.NewCredentials(accessKey, secretKey)
	objectsManager := objects.NewObjectsManager(&objects.ObjectsManagerOptions{
		Options: http_client.Options{Credentials: mac},
	})
	bucket := objectsManager.Bucket(bucketName)

	// 批量操作
	for _, fileName := range fileNames {
		// 图片
		err := bucket.Object("img/" + fileName).SetLifeCycle().DeleteAfterDays(-1).Call(context.Background())
		if err != nil {
			fmt.Printf("七牛SDK出错了%s\n", err)
			return err
		}
	}

	// 封面
	err := bucket.Object("img/" + cover).SetLifeCycle().DeleteAfterDays(-1).Call(context.Background())
	if err != nil {
		fmt.Printf("七牛SDK出错了%s\n", err)
		return err
	}
	return nil
}

// UpdateFileLifecycle 更新单张图片的有效期为永久
func UpdateFileLifecycle(path, fileName, accessKey, secretKey, bucketName string) error {
	// 构建 Bucket 对象
	mac := credentials.NewCredentials(accessKey, secretKey)
	objectsManager := objects.NewObjectsManager(&objects.ObjectsManagerOptions{
		Options: http_client.Options{Credentials: mac},
	})
	bucket := objectsManager.Bucket(bucketName)
	err := bucket.Object(path + "/" + fileName).SetLifeCycle().DeleteAfterDays(-1).Call(context.Background())
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	return nil
}
