package config

// Config 应用配置
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	API      APIConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port string
	Mode string // debug, release
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Path string
}

// APIConfig 外部API配置
type APIConfig struct {
	ShuangseqiuURL string // 双色球API地址
	DaletouURL     string // 大乐透API地址
}

// GetConfig 获取配置
func GetConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port: "8080",
			Mode: "debug",
		},
		Database: DatabaseConfig{
			Path: "./data/lottery.db",
		},
		API: APIConfig{
			// 双色球官方API
			ShuangseqiuURL: "https://www.cwl.gov.cn/cwl_admin/front/cwlkj/search/kjxx/findDrawNotice",
			// 大乐透官方API
			DaletouURL: "https://webapi.sporttery.cn/gateway/lottery/getHistoryPageListV1.qry",
		},
	}
}

