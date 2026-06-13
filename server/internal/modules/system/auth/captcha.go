package auth

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

var captchaStore = base64Captcha.NewMemoryStore(10240, 180)

func GenerateCaptcha() CaptchaResponse {
	driver := base64Captcha.NewDriverDigit(
		global.GVA_CONFIG.Captcha.ImgHeight,
		global.GVA_CONFIG.Captcha.ImgWidth,
		global.GVA_CONFIG.Captcha.KeyLong,
		0.7, 80,
	)
	captcha := base64Captcha.NewCaptcha(driver, captchaStore)
	id, b64s, _, err := captcha.Generate()
	if err != nil {
		global.GVA_LOG.Error("captcha generate failed", zap.Error(err))
		return CaptchaResponse{CaptchaLength: 6}
	}
	openCaptcha := global.GVA_CONFIG.Captcha.OpenCaptcha == 0
	return CaptchaResponse{
		CaptchaId:     id,
		PicPath:       b64s,
		CaptchaLength: global.GVA_CONFIG.Captcha.KeyLong,
		OpenCaptcha:   openCaptcha,
	}
}

func VerifyCaptcha(id, answer string) bool {
	if global.GVA_CONFIG.Captcha.OpenCaptcha != 0 {
		return true
	}
	return captchaStore.Verify(id, answer, true)
}
