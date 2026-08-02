package auth

import (
	"github.com/chengkz2023/My-GVA/server/config"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

var captchaStore = base64Captcha.NewMemoryStore(10240, 180)

func GenerateCaptcha(cfg config.Captcha, log *zap.Logger) CaptchaResponse {
	driver := base64Captcha.NewDriverDigit(
		cfg.ImgHeight,
		cfg.ImgWidth,
		cfg.KeyLong,
		0.7, 80,
	)
	captcha := base64Captcha.NewCaptcha(driver, captchaStore)
	id, b64s, _, err := captcha.Generate()
	if err != nil {
		log.Error("captcha generate failed", zap.Error(err))
		return CaptchaResponse{CaptchaLength: 6}
	}
	openCaptcha := cfg.OpenCaptcha == 0
	return CaptchaResponse{
		CaptchaId:     id,
		PicPath:       b64s,
		CaptchaLength: cfg.KeyLong,
		OpenCaptcha:   openCaptcha,
	}
}

func VerifyCaptcha(id, answer string, openCaptcha int) bool {
	if openCaptcha != 0 {
		return true
	}
	return captchaStore.Verify(id, answer, true)
}
