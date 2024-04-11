package config

import (
	"log"
	"os"
	"reflect"
	"time"

	"DevKit-Neuro-server/internal/validator"
	"go.uber.org/zap/zapcore"

	"github.com/spf13/viper"
)

const (
	DefaultTimeFormat        = "02.01.2006 15:04:05"
	FilterTimeFormat         = "01/02/2006"
	DefaultCompanyUUID       = "d6b1890b-cbb3-4cfd-aba0-e99f814413db"
	DefaultDemoUserLifeTime  = time.Duration(time.Minute * 5)
	FaceRecognizerTimeout    = time.Duration(time.Second * 7)
	ControllerRequestTimeout = time.Duration(time.Second * 15)
)

var (
	DebugLevel = zapcore.DebugLevel
)

// Config используется для хранения конфигурации сервера.
type Config struct {
	Addr                string        `mapstructure:"PUBLIC_ADDR"`
	ThresholdTime       time.Duration `mapstructure:"THRESHOLD_TIME"`
	ThresholdPercentage float32       `mapstructure:"THRESHOLD_PERCENTAGE"`

	NumberChannels int `mapstructure:"NUMBER_CHANNELS"`

	Channel1Key string `mapstructure:"CHANNEL1_KEY"`
	Channel2Key string `mapstructure:"CHANNEL2_KEY"`
	Channel3Key string `mapstructure:"CHANNEL3_KEY"`
	Channel4Key string `mapstructure:"CHANNEL4_KEY"`
	Channel5Key string `mapstructure:"CHANNEL5_KEY"`
	Channel6Key string `mapstructure:"CHANNEL6_KEY"`
	Channel7Key string `mapstructure:"CHANNEL7_KEY"`
	Channel8Key string `mapstructure:"CHANNEL8_KEY"`

	//// MqttBrokerAddr - MQTT broker full addr (пример: tcp://192.168.1.108:1883)
	//MqttBrokerAddr string `mapstructure:"MQTT_BROKER_ADDR" validate:"required"`
	//MqttUsername   string `mapstructure:"MQTT_USERNAME"`
	//MqttPassword   string `mapstructure:"MQTT_PASSWORD"`
	//// MqttKeepAlive - MQTT client keep alive
	//MqttKeepAlive time.Duration `mapstructure:"MQTT_KEEP_ALIVE"`
	//// MqttPingTimeout - MQTT client ping timeout
	//MqttPingTimeout time.Duration `mapstructure:"MQTT_PING_TIMEOUT"`
	//// MqttQOS - MQTT client QoS
	//MqttQOS byte `mapstructure:"MQTT_QOS"`
}

func initDefaultConfig() (v *viper.Viper) {
	v = viper.New()

	v.SetDefault("PUBLIC_ADDR", ":8888")
	v.SetDefault("THRESHOLD_TIME", 30*time.Second)
	v.SetDefault("THRESHOLD_PERCENTAGE", 0.3)

	v.SetDefault("NUMBER_CHANNELS", 8)

	v.SetDefault("CHANNEL1_KEY", "w")
	v.SetDefault("CHANNEL2_KEY", "a")
	v.SetDefault("CHANNEL3_KEY", "s")
	v.SetDefault("CHANNEL4_KEY", "d")
	v.SetDefault("CHANNEL5_KEY", "Space")
	v.SetDefault("CHANNEL6_KEY", "Ctrl")
	v.SetDefault("CHANNEL7_KEY", "shift")
	v.SetDefault("CHANNEL8_KEY", "l")
	return
}

func loadConfigFile(v *viper.Viper, path string) (config Config, err error) {
	v.AddConfigPath(path)
	v.SetConfigName("main")
	v.SetConfigType("env")

	v.AutomaticEnv()

	err = v.ReadInConfig()
	if err != nil {
		return
	}

	configReflectType := reflect.ValueOf(&config).Elem()
	configFieldsCount := configReflectType.NumField()

	err = v.Unmarshal(&config)
	if err != nil {
		return
	}
	time.Now().Unix()

	for fieldInd := 0; fieldInd < configFieldsCount; fieldInd++ {
		configField := configReflectType.Field(fieldInd)

		if configField.Kind() != reflect.Struct {
			continue
		}

		err = v.Unmarshal(configField.Addr().Interface())
		if err != nil {
			return
		}
	}

	return
}

func loadConfigEnv(v *viper.Viper) (config Config, err error) {
	envNameList := envNameListByConfig(reflect.TypeOf(config))
	for _, envName := range envNameList {
		err = v.BindEnv(envName, envName)
		if err != nil {
			return
		}
	}

	err = v.Unmarshal(&config)
	return
}

func envNameListByConfig(configType reflect.Type) (envNameList []string) {
	configFieldsCount := configType.NumField()
	envNameList = make([]string, 0, configFieldsCount)

	for fieldInd := 0; fieldInd < configFieldsCount; fieldInd++ {
		configField := configType.Field(fieldInd)

		if configField.Type.Kind() == reflect.Struct {
			envNameList = append(envNameList, envNameListByConfig(configField.Type)...)
		}

		envNameList = append(envNameList, configField.Tag.Get("mapstructure"))
	}
	return
}

func LoadConfig(appValidator *validator.AppValidator) (config Config, err error) {
	v := initDefaultConfig()

	if _, err = os.Stat("./main.env"); err == nil {
		config, err = loadConfigFile(v, "./")
	} else {
		log.Println("Loading config from env...")
		config, err = loadConfigEnv(v)
	}

	if err = appValidator.Struct(&config); err != nil {
		err = appValidator.ErrorTranslated(err)
		return
	}

	return config, err
}
