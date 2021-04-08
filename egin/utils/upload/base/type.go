package base

type Interface interface {
	Save(bucketName string, localFile string, saveFile string) (string, error)
}
