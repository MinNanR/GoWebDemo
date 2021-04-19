package main

import (
	"net/http"
	"strings"
)

const (
	TABLET_TYPE = "tablet"
	MOBILE_TYPE = "mobile"
	NORMAL_TYPE = "normal"
)

const (
	ANDROID_PLATFORM = "android"
	IOS_PLATFORM     = "ios"
	UNKNOWN_PLATFORM = "unknown"
)

var mobileUserAgentKeywords = []string{
	"blackberry", "webos", "ipod", "lge vx", "midp", "maemo", "mmp", "mobile",
	"netfront", "hiptop", "nintendo DS", "novarra", "openweb", "opera mobi",
	"opera mini", "palm", "psp", "phone", "smartphone", "symbian", "up.browser",
	"up.link", "wap", "windows ce",
}

var mobileUserAgentPrefixes = []string{
	"w3c ", "w3c-", "acs-", "alav", "alca", "amoi", "audi", "avan", "benq",
	"bird", "blac", "blaz", "brew", "cell", "cldc", "cmd-", "dang", "doco",
	"eric", "hipt", "htc_", "inno", "ipaq", "ipod", "jigs", "kddi", "keji",
	"leno", "lg-c", "lg-d", "lg-g", "lge-", "lg/u", "maui", "maxo", "midp",
	"mits", "mmef", "mobi", "mot-", "moto", "mwbp", "nec-", "newt", "noki",
	"palm", "pana", "pant", "phil", "play", "port", "prox", "qwap", "sage",
	"sams", "sany", "sch-", "sec-", "send", "seri", "sgh-", "shar", "sie-",
	"siem", "smal", "smar", "sony", "sph-", "symb", "t-mo", "teli", "tim-",
	"tosh", "tsm-", "upg1", "upsi", "vk-v", "voda", "wap-", "wapa", "wapi",
	"wapp", "wapr", "webc", "winw", "winw", "xda ", "xda-",
}

var tabletUserAgentKeywords = []string{
	"ipad", "playbook", "hp-tablet", "kindle",
}

func getDevice(request http.Request) (string, string) {
	header := request.Header
	userAgent := header.Get("User-Agent")
	if len(userAgent) > 0 {
		userAgent = strings.ToLower(userAgent)
		//安卓平板
		if strings.Contains(userAgent, "android") && !strings.Contains(userAgent, "mobile") {
			return TABLET_TYPE, ANDROID_PLATFORM
		}
		//ipad
		if strings.Contains(userAgent, "ipad") {
			return TABLET_TYPE, IOS_PLATFORM
		}
		//kindle
		if strings.Contains(userAgent, "silk") && !strings.Contains(userAgent, "mobile") {
			return TABLET_TYPE, UNKNOWN_PLATFORM
		}
		for _, keyword := range tabletUserAgentKeywords {
			if strings.Contains(userAgent, keyword) {
				return TABLET_TYPE, UNKNOWN_PLATFORM
			}
		}
	}

	if len(header.Get("x-wap-profile")) > 0 || len(header.Get("Profile")) > 0 {
		if len(userAgent) > 0 {
			//安卓手机
			if strings.Contains(userAgent, "android") {
				return MOBILE_TYPE, ANDROID_PLATFORM
			}
			//IOS设备
			if strings.Contains(userAgent, "iphone") || strings.Contains(userAgent,
				"ipod") || strings.Contains(userAgent, "ipad") {
				return MOBILE_TYPE, IOS_PLATFORM
			}
		}
		return MOBILE_TYPE, UNKNOWN_PLATFORM
	}

	if len(userAgent) > 4 {
		prefix := strings.ToLower(userAgent[0:4])
		for _, keyword := range mobileUserAgentPrefixes {
			if prefix == keyword {
				return MOBILE_TYPE, UNKNOWN_PLATFORM
			}
		}
	}

	accept := header.Get("Accept")
	if len(accept) > 0 {
		return MOBILE_TYPE, UNKNOWN_PLATFORM
	}

	if len(userAgent) > 0 {
		if strings.Contains(userAgent, "android") {
			return MOBILE_TYPE, ANDROID_PLATFORM
		}
		if strings.Contains(userAgent, "iphone") || strings.Contains(userAgent, "ipod") || strings.Contains(userAgent,
			"ipad") {
			return MOBILE_TYPE, IOS_PLATFORM
		}
		for _, keyword := range mobileUserAgentKeywords {
			if strings.Contains(userAgent, keyword) {
				return MOBILE_TYPE, UNKNOWN_PLATFORM
			}
		}
	}

	return NORMAL_TYPE, UNKNOWN_PLATFORM
}
