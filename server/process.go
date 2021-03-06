package server

import (
	"captcha-zh/config"
	"captcha-zh/captcha"

	"log"
	"time"
)

func captchaGenerate(size int) []string {
	s := make([]string, 0)
	for i := 0; i < size; i++ {
		Topic := captcha.RandTopic()
		fileName := captcha.RandomName() + ".gif"
		captcha.Draw(Topic.Subject, config.TConfig.Paths.Path + config.PATH_CONFIG_IMAGE_TEMP + fileName)
		s = append(s, fileName + config.SEPARATOR_VERTICAL_LINE + Topic.Result)
	}
	return s
}

func Start() {
	c := config.TConfig.CaptchaSys
	// 初始化生成
	captchas := captchaGenerate(c.Initial_count)

	// 更新到容器
	captcha.CaptchaContainer.Append(captchas...)
	log.Print("init success.")

	// 定时检查
	ticker := time.NewTicker(time.Second * time.Duration(c.Check_interval))
	go func() {
		for _ = range ticker.C {
			go workder()
		}
	}()
}

func workder() {
	if !captcha.CaptchaContainer.UpdateNeed() {
		return
	}

	// 新的
	captchas := captchaGenerate(config.TConfig.CaptchaSys.Update_count)

	captcha.CaptchaContainer.Lock()
	captcha.CaptchaContainer.Update(captchas...)
	captcha.CaptchaContainer.Unlock()

	log.Print("update suceess.")
}
