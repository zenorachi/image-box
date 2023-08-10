package logger

import "github.com/sirupsen/logrus"

const (
	SignUpHandler   = "signUp"
	SignInHandler   = "signIn"
	RefreshHandler  = "refreshHandler"
	UploadHandler   = "uploadHandler"
	GetFilesHandler = "getFilesHandler"
	AuthMiddleware  = "authMiddleware"
	FilesMiddleware = "filesMiddleware"
)

func logFields(handler string) logrus.Fields {
	return logrus.Fields{"handler": handler}
}

func LogError(handler string, err error) {
	logrus.WithFields(logFields(handler)).Error(err)
}
