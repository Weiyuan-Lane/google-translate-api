package server

import (
	"strconv"

	"github.com/weiyuan-lane/google-translate-api/internal/services/googletranslatewrapper"
	httptransport "github.com/weiyuan-lane/google-translate-api/internal/transports/http"
	"github.com/weiyuan-lane/google-translate-api/internal/utils/config"
	"github.com/weiyuan-lane/google-translate-api/internal/utils/googletranslate"
	loggerutils "github.com/weiyuan-lane/google-translate-api/internal/utils/logger"
)

func Init() {
	appConfig := config.ApplicationConfig()

	logger := loggerutils.New(
		appConfig.AppName,
		appConfig.IsDevEnv,
	)

	googleTranslateV2Client := googletranslate.InitTranslateV2Client(
		appConfig.GoogleTranslateV2APIKey,
	)
	defer googleTranslateV2Client.Close()

	googleTranslateV3Client := googletranslate.InitTranslateV3Client()
	defer googleTranslateV2Client.Close()

	translateV3Wrapper := googletranslatewrapper.NewTranslateV3Wrapper(
		googleTranslateV3Client,
		appConfig.GoogleTranslateV3ProjectKey,
	)
	translateV2Wrapper := googletranslatewrapper.NewTranslateV2WrapperWithV3Wrapper(
		googleTranslateV2Client,
		translateV3Wrapper,
	)

	httpServer := httptransport.HttpServer{
		LivelinessProbePort:      strconv.Itoa(appConfig.LivenessPort),
		Port:                     strconv.Itoa(appConfig.Port),
		Logger:                   logger,
		GracefulShutdownSeconds:  appConfig.GracefulShutdownSeconds,
		EnableHTTP2:              appConfig.EnableHTTP2,
		GoogleTranslateV2Wrapper: translateV2Wrapper,
		GoogleTranslateV3Wrapper: translateV3Wrapper,
	}

	httpServer.ListenAndServe()
}
