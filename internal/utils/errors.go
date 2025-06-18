package utils

import (
	"fmt"
	"os"
)

// 에러 타입 정의
type ErrorType string

const (
	ErrorTypeValidation ErrorType = "VALIDATION"
	ErrorTypeDatabase   ErrorType = "DATABASE"
	ErrorTypeUserInput  ErrorType = "USER_INPUT"
	ErrorTypeSystem     ErrorType = "SYSTEM"
	ErrorTypeNotFound   ErrorType = "NOT_FOUND"
)

// 애플리케이션 에러 구조체
type AppError struct {
	Type    ErrorType // 에러 종류
	Message string    // 에러 메시지
	Err     error
	Code    int // 종료 코드
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

// AppError 생성 함수
func NewError(errType ErrorType, message string, err error) *AppError {
	return &AppError{
		Type:    errType,
		Message: message,
		Err:     err,
		Code:    getErrorCode(errType),
	}
}

// 에러 타입별 종료 코드 반환
func getErrorCode(errType ErrorType) int {
	switch errType {
	case ErrorTypeValidation:
		return 1
	case ErrorTypeUserInput:
		return 2
	case ErrorTypeNotFound:
		return 3
	case ErrorTypeDatabase:
		return 5
	case ErrorTypeSystem:
		return 6
	default:
		return 1
	}
}

// 에러 출력하고 종료
func HandleError(err error, message string) {
	if err == nil {
		return
	}

	var appErr *AppError
	if e, ok := err.(*AppError); ok {
		appErr = e
	} else {
		appErr = NewError(ErrorTypeSystem, message, err)
	}

	emoji := getErrorEmoji(appErr.Type)

	fmt.Fprintf(os.Stderr, "%s %s\n", emoji, appErr.Message)
	Error("%s %s", emoji, appErr.Message)

	os.Exit(appErr.Code)
}

// 에러 타입별 이모지 반환
func getErrorEmoji(errType ErrorType) string {
	switch errType {
	case ErrorTypeValidation:
		return "⚠️"
	case ErrorTypeUserInput:
		return "❌"
	case ErrorTypeNotFound:
		return "🔍"
	case ErrorTypeDatabase:
		return "🗄️"
	case ErrorTypeSystem:
		return "💥"
	default:
		return "❗"
	}
}

// 유효성 검사 에러
func ValidationError(message string, err error) *AppError {
	return NewError(ErrorTypeValidation, message, err)
}

// DB 관련 에러
func DatabaseError(message string, err error) *AppError {
	return NewError(ErrorTypeDatabase, message, err)
}

// 사용자 입력 에러
func UserInputError(message string, err error) *AppError {
	return NewError(ErrorTypeUserInput, message, err)
}

// 데이터 없음 에러
func NotFoundError(message string, err error) *AppError {
	return NewError(ErrorTypeNotFound, message, err)
}

// 시스템 내부 에러
func SystemError(message string, err error) *AppError {
	return NewError(ErrorTypeSystem, message, err)
}
