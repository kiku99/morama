package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/kiku99/morama/internal/config"
)

// 로그 레벨 정의
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarning
	LogLevelError
)

// Logger 구조체: 레벨별 로거, 로그 파일, 현재 레벨 포함
type Logger struct {
	debug   *log.Logger
	info    *log.Logger
	warning *log.Logger
	error   *log.Logger
	file    *os.File
	level   LogLevel
}

// Logger 인스턴스 생성
func NewLogger(level LogLevel) (*Logger, error) {
	logDir, err := getLogDir()
	if err != nil {
		return nil, err
	}

	// 오래된 로그 삭제 (7일 보관)
	_ = cleanupOldLogs(logDir, 7)

	// 로그 디렉토리 없으면 생성
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return nil, err
	}

	// 로그 파일 경로 설정
	timestamp := time.Now().Format("2006-01-02")
	logFile := filepath.Join(logDir, fmt.Sprintf("morama-%s.log", timestamp))

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	flags := log.Ldate | log.Ltime | log.Lshortfile

	return &Logger{
		debug:   log.New(file, "[DEBUG] ", flags),
		info:    log.New(file, "[INFO] ", flags),
		warning: log.New(file, "[WARN] ", flags),
		error:   log.New(file, "[ERROR] ", flags),
		file:    file,
		level:   level,
	}, nil
}

// 로그 디렉토리 경로 반환 (예: ~/.morama/logs)
func getLogDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".morama", "logs"), nil
}

// 오래된 로그 삭제 (keepDays일 이상된 파일 삭제)
func cleanupOldLogs(dir string, keepDays int) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	cutoff := time.Now().AddDate(0, 0, -keepDays)
	for _, file := range files {
		info, err := file.Info()
		if err != nil {
			continue
		}
		if info.ModTime().Before(cutoff) {
			_ = os.Remove(filepath.Join(dir, file.Name()))
		}
	}
	return nil
}

// 디버그 로그 출력
func (l *Logger) Debug(format string, v ...interface{}) {
	if l.level <= LogLevelDebug {
		l.debug.Printf(format, v...)
	}
}

// 정보 로그 출력
func (l *Logger) Info(format string, v ...interface{}) {
	if l.level <= LogLevelInfo {
		l.info.Printf(format, v...)
	}
}

// 경고 로그 출력
func (l *Logger) Warning(format string, v ...interface{}) {
	if l.level <= LogLevelWarning {
		l.warning.Printf(format, v...)
	}
}

// 에러 로그 출력
func (l *Logger) Error(format string, v ...interface{}) {
	if l.level <= LogLevelError {
		l.error.Printf(format, v...)
	}
}

// 로그 파일 닫기
func (l *Logger) Close() error {
	if l.file != nil {
		return l.file.Close()
	}
	return nil
}

// 사용자 행동 로그 기록
func (l *Logger) LogUserAction(action, details string) {
	l.Info("USER_ACTION: %s - %s", action, details)
}

// 커맨드 실행 로그 기록
func (l *Logger) LogCommandExecution(command string, args []string, duration time.Duration) {
	l.Info("COMMAND: %s %v (duration: %v)", command, args, duration)
}

// DB 작업 로그 기록
func (l *Logger) LogDatabaseOperation(operation, table string, duration time.Duration) {
	l.Debug("DB_OP: %s on %s (duration: %v)", operation, table, duration)
}

// 전역 Logger 인스턴스
var globalLogger *Logger

// 전역 Logger 초기화 (직접 레벨 지정)
func InitLogger(level LogLevel) error {
	var err error
	globalLogger, err = NewLogger(level)
	return err
}

// 전역 Logger 초기화 (config.yaml 기반)
func InitLoggerFromConfig() error {
	cfg := config.GetConfig()
	level := LogLevelInfo
	if cfg.DebugMode {
		level = LogLevelDebug
	} else {
		level = LogLevelWarning
	}
	return InitLogger(level)
}

// 전역 Logger 반환
func GetLogger() *Logger {
	if globalLogger == nil {
		InitLogger(LogLevelInfo)
	}
	return globalLogger
}

// 전역 디버그 로그 출력
func Debug(format string, v ...interface{}) {
	GetLogger().Debug(format, v...)
}

// 전역 정보 로그 출력
func Info(format string, v ...interface{}) {
	GetLogger().Info(format, v...)
}

// 전역 경고 로그 출력
func Warning(format string, v ...interface{}) {
	GetLogger().Warning(format, v...)
}

// 전역 에러 로그 출력
func Error(format string, v ...interface{}) {
	GetLogger().Error(format, v...)
}

// 사용자 행동 기록
func LogUserAction(action, details string) {
	GetLogger().LogUserAction(action, details)
}

// 명령 실행 로그
func LogCommandExecution(command string, args []string, duration time.Duration) {
	GetLogger().LogCommandExecution(command, args, duration)
}

// DB 작업 로그
func LogDatabaseOperation(operation, table string, duration time.Duration) {
	GetLogger().LogDatabaseOperation(operation, table, duration)
}
