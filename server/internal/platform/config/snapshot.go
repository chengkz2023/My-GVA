package config

type Snapshot struct {
	System  SystemSnapshot  `json:"system"`
	JWT     JWTSnapshot     `json:"jwt"`
	MySQL   MySQLSnapshot   `json:"mysql"`
	Redis   RedisSnapshot   `json:"redis"`
	Zap     ZapSnapshot     `json:"zap"`
	CORS    CORSSnapshot    `json:"cors"`
	Local   LocalSnapshot   `json:"local"`
	Captcha CaptchaSnapshot `json:"captcha"`
}

type SystemSnapshot struct {
	RouterPrefix       string `json:"routerPrefix"`
	Addr               int    `json:"addr"`
	LimitCountIP       int    `json:"limitCountIp"`
	LimitTimeIP        int    `json:"limitTimeIp"`
	UseStrictAuth      bool   `json:"useStrictAuth"`
	DisableAutoMigrate bool   `json:"disableAutoMigrate"`
}

type JWTSnapshot struct {
	ExpiresTime   string `json:"expiresTime"`
	BufferTime    string `json:"bufferTime"`
	Issuer        string `json:"issuer"`
	SigningKeySet bool   `json:"signingKeySet"`
}

type MySQLSnapshot struct {
	Path         string `json:"path"`
	Port         string `json:"port"`
	DBName       string `json:"dbName"`
	Username     string `json:"username"`
	Config       string `json:"config"`
	MaxIdleConns int    `json:"maxIdleConns"`
	MaxOpenConns int    `json:"maxOpenConns"`
	LogMode      string `json:"logMode"`
	LogZap       bool   `json:"logZap"`
	Prefix       string `json:"prefix"`
	Singular     bool   `json:"singular"`
	Engine       string `json:"engine"`
	PasswordSet  bool   `json:"passwordSet"`
}

type RedisSnapshot struct {
	Enable       bool     `json:"enable"`
	Name         string   `json:"name"`
	Addr         string   `json:"addr"`
	DB           int      `json:"db"`
	UseCluster   bool     `json:"useCluster"`
	ClusterAddrs []string `json:"clusterAddrs"`
	PasswordSet  bool     `json:"passwordSet"`
}

type ZapSnapshot struct {
	Level         string `json:"level"`
	Prefix        string `json:"prefix"`
	Format        string `json:"format"`
	Director      string `json:"director"`
	EncodeLevel   string `json:"encodeLevel"`
	StacktraceKey string `json:"stacktraceKey"`
	ShowLine      bool   `json:"showLine"`
	LogInConsole  bool   `json:"logInConsole"`
	RetentionDay  int    `json:"retentionDay"`
}

type CORSSnapshot struct {
	Mode           string `json:"mode"`
	WhitelistCount int    `json:"whitelistCount"`
}

type LocalSnapshot struct {
	Path      string `json:"path"`
	StorePath string `json:"storePath"`
}

type CaptchaSnapshot struct {
	KeyLong            int `json:"keyLong"`
	ImgWidth           int `json:"imgWidth"`
	ImgHeight          int `json:"imgHeight"`
	OpenCaptcha        int `json:"openCaptcha"`
	OpenCaptchaTimeout int `json:"openCaptchaTimeout"`
}

func SafeSnapshot(cfg Config) Snapshot {
	return Snapshot{
		System: SystemSnapshot{
			RouterPrefix:       cfg.System.RouterPrefix,
			Addr:               cfg.System.Addr,
			LimitCountIP:       cfg.System.LimitCountIP,
			LimitTimeIP:        cfg.System.LimitTimeIP,
			UseStrictAuth:      cfg.System.UseStrictAuth,
			DisableAutoMigrate: cfg.System.DisableAutoMigrate,
		},
		JWT: JWTSnapshot{
			ExpiresTime:   cfg.JWT.ExpiresTime,
			BufferTime:    cfg.JWT.BufferTime,
			Issuer:        cfg.JWT.Issuer,
			SigningKeySet: cfg.JWT.SigningKey != "",
		},
		MySQL: MySQLSnapshot{
			Path:         cfg.Mysql.Path,
			Port:         cfg.Mysql.Port,
			DBName:       cfg.Mysql.Dbname,
			Username:     cfg.Mysql.Username,
			Config:       cfg.Mysql.Config,
			MaxIdleConns: cfg.Mysql.MaxIdleConns,
			MaxOpenConns: cfg.Mysql.MaxOpenConns,
			LogMode:      cfg.Mysql.LogMode,
			LogZap:       cfg.Mysql.LogZap,
			Prefix:       cfg.Mysql.Prefix,
			Singular:     cfg.Mysql.Singular,
			Engine:       cfg.Mysql.Engine,
			PasswordSet:  cfg.Mysql.Password != "",
		},
		Redis: RedisSnapshot{
			Enable:       cfg.Redis.Enable,
			Name:         cfg.Redis.Name,
			Addr:         cfg.Redis.Addr,
			DB:           cfg.Redis.DB,
			UseCluster:   cfg.Redis.UseCluster,
			ClusterAddrs: append([]string(nil), cfg.Redis.ClusterAddrs...),
			PasswordSet:  cfg.Redis.Password != "",
		},
		Zap: ZapSnapshot{
			Level:         cfg.Zap.Level,
			Prefix:        cfg.Zap.Prefix,
			Format:        cfg.Zap.Format,
			Director:      cfg.Zap.Director,
			EncodeLevel:   cfg.Zap.EncodeLevel,
			StacktraceKey: cfg.Zap.StacktraceKey,
			ShowLine:      cfg.Zap.ShowLine,
			LogInConsole:  cfg.Zap.LogInConsole,
			RetentionDay:  cfg.Zap.RetentionDay,
		},
		CORS: CORSSnapshot{
			Mode:           cfg.Cors.Mode,
			WhitelistCount: len(cfg.Cors.Whitelist),
		},
		Local: LocalSnapshot{
			Path:      cfg.Local.Path,
			StorePath: cfg.Local.StorePath,
		},
		Captcha: CaptchaSnapshot{
			KeyLong:            cfg.Captcha.KeyLong,
			ImgWidth:           cfg.Captcha.ImgWidth,
			ImgHeight:          cfg.Captcha.ImgHeight,
			OpenCaptcha:        cfg.Captcha.OpenCaptcha,
			OpenCaptchaTimeout: cfg.Captcha.OpenCaptchaTimeOut,
		},
	}
}
