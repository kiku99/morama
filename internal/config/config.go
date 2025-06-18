package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config 설정 구조체
type Config struct {
	Display   DisplayConfig `yaml:"display"`
	Search    SearchConfig  `yaml:"search"`
	DebugMode bool          `yaml:"debug_mode"`
}

// DisplayConfig 출력 관련 설정
type DisplayConfig struct {
	DateFormat  string  `yaml:"date_format"`  // 날짜 출력 형식
	RatingScale float64 `yaml:"rating_scale"` // 평점의 최대값
	ShowEmojis  bool    `yaml:"show_emojis"`  // 출력에 이모지를 보여줄지
}

// SearchConfig 검색 관련 설정
type SearchConfig struct {
	FuzzyMatch    bool `yaml:"fuzzy_match"`
	CaseSensitive bool `yaml:"case_sensitive"`
	MaxResults    int  `yaml:"max_results"`
}

// DefaultConfig 기본 설정값 반환
func DefaultConfig() *Config {
	return &Config{
		Display: DisplayConfig{
			DateFormat:  "2006-01-02",
			RatingScale: 5.0,
			ShowEmojis:  true,
		},
		Search: SearchConfig{
			FuzzyMatch:    true,
			CaseSensitive: false,
			MaxResults:    50,
		},
		DebugMode: false,
	}
}

// LoadConfig 설정 파일 불러오기 (없으면 기본값 생성)
func LoadConfig() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config := DefaultConfig()
		if err := SaveConfig(config); err != nil {
			return nil, err
		}
		return config, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return mergeWithDefaults(config), nil
}

// SaveConfig 설정 파일 저장
func SaveConfig(config *Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(configPath, data, 0644)
}

// 설정 파일 경로 반환
func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".morama", "config.yaml"), nil
}

// 누락된 필드에 기본값 채워 넣기
func mergeWithDefaults(config Config) *Config {
	defaults := DefaultConfig()

	if config.Display.DateFormat == "" {
		config.Display.DateFormat = defaults.Display.DateFormat
	}
	if config.Display.RatingScale == 0 {
		config.Display.RatingScale = defaults.Display.RatingScale
	}
	if config.Search.MaxResults == 0 {
		config.Search.MaxResults = defaults.Search.MaxResults
	}

	config.DebugMode = config.DebugMode || defaults.DebugMode

	return &config
}

// 전역 설정 인스턴스 반환
var globalConfig *Config

func GetConfig() *Config {
	if globalConfig == nil {
		var err error
		globalConfig, err = LoadConfig()
		if err != nil {
			globalConfig = DefaultConfig()
		}
	}
	return globalConfig
}
