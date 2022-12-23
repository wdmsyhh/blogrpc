package util

import (
	"fmt"
	"os"
	"strings"

	conf "github.com/spf13/viper"
)

const (
	TENCENT_CAPTCHA_SCENE_ACCOUNT  = "account"
	TENCENT_CAPTCHA_SCENE_SMS      = "sms"
	TENCENT_CAPTCHA_SCENE_CAMPAIGN = "campaign"
	TENCENT_CAPTCHA_SCENE_OTHER    = "other"
)

func GetAliyunAKId() string {
	return os.Getenv("ALIYUN_ACCESS_KEY_ID")
}

func GetAliyunAKSecret() string {
	return os.Getenv("ALIYUN_ACCESS_KEY_SECRET")
}

func GetAliyunGreenWebAKId() string {
	k := os.Getenv("ALIYUN_GREENWEB_ACCESS_KEY_ID")
	if k == "" {
		return GetAliyunAKId()
	}
	return k
}

func GetAliyunGreenWebAKSecret() string {
	k := os.Getenv("ALIYUN_GREENWEB_ACCESS_KEY_SECRET")
	if k == "" {
		return GetAliyunAKSecret()
	}
	return k
}

func GetAliyunMqAKId() string {
	k := os.Getenv("ALIYUN_MQ_ACCESS_KEY_ID")
	if k == "" {
		return GetAliyunAKId()
	}
	return k
}

func GetAliyunMqAKSecret() string {
	k := os.Getenv("ALIYUN_MQ_ACCESS_KEY_SECRET")
	if k == "" {
		return GetAliyunAKSecret()
	}
	return k
}

func GetAliyunSmsAKId() string {
	k := os.Getenv("ALIYUN_SMS_ACCESS_KEY_ID")
	if k == "" {
		return GetAliyunAKId()
	}
	return k
}

func GetAliyunSmsAKSecret() string {
	k := os.Getenv("ALIYUN_SMS_ACCESS_KEY_SECRET")
	if k == "" {
		return GetAliyunAKSecret()
	}
	return k
}

func GetTianyuSecretId() string {
	return os.Getenv("TIANYU_SECRET_ID")
}

func GetTianyuSecretKey() string {
	return os.Getenv("TIANYU_SECRET_KEY")
}

func GetAliyunSMTPPassword() string {
	return os.Getenv("SMTP_PASSWORD")
}

func GetBackendDomain() string {
	return conf.GetString("redirect-backend-url")
}

func GetTencentLbsKey() string {
	return os.Getenv("TENCENT_LBS_KEY")
}

func GetTencentMapApi() string {
	return os.Getenv("TENCENT_MAP_API")
}

func GetYouZanYunClientId() string {
	return os.Getenv("YOUZANYUN_CLIENT_ID")
}

func GetYouZanYunClientSecret() string {
	return os.Getenv("YOUZANYUN_CLIENT_SECRET")
}

func GetTaobaoAppKey() string {
	return os.Getenv("TAOBAO_APP_KEY")
}

func GetTaobaoAppSecret() string {
	return os.Getenv("TAOBAO_APP_SECRET")
}

func GetTaobaoCrmAppKey() string {
	return os.Getenv("TAOBAO_CRM_APP_KEY")
}

func GetTaobaoCrmAppSecret() string {
	return os.Getenv("TAOBAO_CRM_APP_SECRET")
}

func GetKDNiaoBusinessId() string {
	return os.Getenv("KDNIAO_BUSINESS_ID")
}

func GetKDNiaoAppKey() string {
	return os.Getenv("KDNIAO_APP_KEY")
}

func GetMqProvider() string {
	return os.Getenv("MQ_PROVIDER")
}

func GetMqAKId() string {
	k := os.Getenv("MQ_ACCESS_KEY_ID")
	if k == "" {
		return "RocketMQ"
	}
	return k
}

func GetMqNamespace() string {
	return os.Getenv("MQ_NAMESPACE")
}

func GetMqAKSecret() string {
	k := os.Getenv("MQ_ACCESS_KEY_SECRET")
	if k == "" {
		return "12345678"
	}
	return k
}

func GetTencentCaptchaAppID(scene string) string {
	appId, _ := getTencentCaptchaConfigs(scene)
	return appId
}

func GetTencentCaptchaAppSecret(scene string) string {
	_, appSecret := getTencentCaptchaConfigs(scene)
	return appSecret
}

func GetFrontendDomains() []string {
	return strings.Split(os.Getenv("FRONTEND_DOMAINS"), " ")
}

func GetPortalH5V2Url() string {
	return os.Getenv("PORTAL_H5V2_URL")
}

func getTencentCaptchaConfigs(scene string) (string, string) {
	isValidScene := strContains([]string{
		TENCENT_CAPTCHA_SCENE_ACCOUNT,
		TENCENT_CAPTCHA_SCENE_SMS,
		TENCENT_CAPTCHA_SCENE_CAMPAIGN,
		TENCENT_CAPTCHA_SCENE_OTHER,
	}, scene)
	if !isValidScene {
		return "", ""
	}

	if conf.GetString("env") == "local" {
		return getTencentCaptchaFromConf(scene)
	}

	return getTencentCaptchaFromEnv(scene)
}

func getTencentCaptchaFromConf(scene string) (string, string) {
	captchaConf := conf.GetString("tencent-captcha-config")
	return parseTencentCaptchaConf(captchaConf, scene)
}

func getTencentCaptchaFromEnv(scene string) (string, string) {
	captchaConf := os.Getenv("TENCENT_CAPTCHA_CONFIGS")
	return parseTencentCaptchaConf(captchaConf, scene)
}

func parseTencentCaptchaConf(conf, scene string) (string, string) {
	for _, config := range strings.Split(conf, ",") {
		config = strings.TrimSpace(config)
		parts := strings.Split(config, " ")

		if len(parts) < 3 {
			continue
		}

		if parts[0] != scene {
			continue
		}

		return parts[1], parts[2]
	}

	return "", ""
}

func GetSreadminJobName() string {
	return os.Getenv("JOB_NAME")
}

func GetCacheHost() string {
	return os.Getenv("CACHE_HOST")
}

func GetCachePort() string {
	return os.Getenv("CACHE_PORT")
}

func GetCachePassword() string {
	return os.Getenv("CACHE_PASSWORD")
}

func GetResqueHost() string {
	return os.Getenv("RESQUE_HOST")
}

func GetResquePort() string {
	return os.Getenv("RESQUE_PORT")
}

func GetResquePassword() string {
	return os.Getenv("RESQUE_PASSWORD")
}

func GetMongoMasterDsn() string {
	return os.Getenv("MONGO_MASTER_DSN")
}

func GetMongoMasterReplset() string {
	return os.Getenv("MONGO_MASTER_REPLSET")
}

func GetElasticsearchUrl() string {
	return os.Getenv("ELASTICSEARCH_URL")
}

func GetElasticsearchUsername() string {
	return os.Getenv("ELASTICSEARCH_USERNAME")
}

func GetElasticsearchPassword() string {
	return os.Getenv("ELASTICSEARCH_PASSWORD")
}

func GetWechatAppId() string {
	return os.Getenv("WECHAT_PROVIDER_APP_ID")
}

func GetByteDanceAppId() string {
	return os.Getenv("BYTEDANCE_PROVIDER_APP_ID")
}

func GetWeiboAppKey() string {
	return os.Getenv("WEIBO_APP_KEY")
}

func GetWeiboAppSecret() string {
	return os.Getenv("WEIBO_APP_SECRET")
}

func GetFeieUser() string {
	return os.Getenv("FEIE_USER")
}

func GetFeieUkey() string {
	return os.Getenv("FEIE_UKEY")
}

func GetMongoAppName() string {
	return fmt.Sprintf("%s.%s", os.Getenv("K8S_SERVICE_NAME"), os.Getenv("K8S_SERVICE_NAMESPACE"))
}

func GetOSSProvider() string {
	k := os.Getenv("OSS_PROVIDER")
	if k == "" {
		return "aliyun"
	}
	return k
}

func GetOSSAccessKeyId() string {
	k := os.Getenv("OSS_ACCESS_KEY_ID")
	if k == "" {
		return GetAliyunAKId()
	}
	return k
}

func GetOSSAccessKeySecret() string {
	k := os.Getenv("OSS_ACCESS_KEY_SECRET")
	if k == "" {
		return GetAliyunAKSecret()
	}
	return k
}

func DisableMarketoProcesser() bool {
	disable := os.Getenv("DISABLE_MARKETO_PROCESSER")
	return disable == "true"
}

func GetWechatProviderAppSecret() string {
	return os.Getenv("WECHAT_PROVIDER_APP_SECRET")
}

func UsePlatformProxyWechatProviderApp() bool {
	usePlatform := os.Getenv("PLATFORM_PROXY_WECHAT_PROVIDER")
	return usePlatform == "true"
}

func GetWechatProviderToken() string {
	return os.Getenv("WECHAT_PROVIDER_TOKEN")
}

func GetWechatProviderEncodingAESKey() string {
	return os.Getenv("WECHAT_PROVIDER_ENCODING_AES_KEY")
}

func GetRegistryHost() string {
	k := os.Getenv("REGISTRY_HOST")
	if k == "" {
		return "cr.maiscrm.com"
	}
	return k
}

func GetYiwiseAppKey() string {
	return os.Getenv("YIWISE_APP_KEY")
}

func GetYiwiseAppSecret() string {
	return os.Getenv("YIWISE_APP_SECRET")
}

func GetYiwiseTenantSign() string {
	return os.Getenv("YIWISE_TENANT_SIGN")
}
