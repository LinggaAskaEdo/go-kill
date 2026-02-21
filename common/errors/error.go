package errors

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/linggaaskaedo/go-kill/common/preference"
	"github.com/palantir/stacktrace"
)

type AppError struct {
	Code       Code    `json:"code"`
	Message    string  `json:"message"`
	DebugError *string `json:"debug,omitempty"`
	sys        error
}

func init() {
	svcError = map[ServiceType]ErrorMessage{
		COMMON: ErrorMessages,
	}
}

func Compile(service ServiceType, err error, lang string, debugMode bool) (int, AppError) {
	// Capture debug error string if requested
	var debugErr *string
	if debugMode {
		if errStr := err.Error(); errStr != "" {
			debugErr = &errStr
		}
	}

	code := ErrCode(err)

	// 1. Check common errors first
	if status, msg, ok := getCommonErrorMessage(code, lang); ok {
		return status, AppError{
			Code:       code,
			Message:    msg,
			sys:        err,
			DebugError: debugErr,
		}
	}

	// 2. Verify service exists in error map
	serviceMsgs, exists := svcError[service]
	if !exists {
		return http.StatusInternalServerError, AppError{
			Code:       code,
			Message:    "service error not defined!",
			sys:        err,
			DebugError: debugErr,
		}
	}

	// 3. Verify error code exists for this service
	errMsg, ok := serviceMsgs[code]
	if !ok {
		return http.StatusInternalServerError, AppError{
			Code:       code,
			Message:    "error message not defined!",
			sys:        err,
			DebugError: debugErr,
		}
	}

	// 4. Select message based on language
	msg := selectLanguageMessage(errMsg, lang)

	// 5. Apply annotation formatting if needed
	if errMsg.HasAnnotation {
		msg = formatAnnotatedMessage(msg, err.Error())
	}

	// 6. Special handling for HTTP validator errors
	if code == CodeHTTPValidatorError && err.Error() != "" {
		msg = strings.Split(err.Error(), "\n ---")[0]
	}

	return errMsg.StatusCode, AppError{
		Code:       code,
		Message:    msg,
		sys:        err,
		DebugError: debugErr,
	}
}

// getCommonErrorMessage looks up a common error message for the given code.
// Returns (statusCode, message, true) if found, otherwise (0, "", false).
func getCommonErrorMessage(code stacktrace.ErrorCode, lang string) (int, string, bool) {
	errMsg, ok := svcError[COMMON][code]
	if !ok {
		return 0, "", false
	}
	return errMsg.StatusCode, selectLanguageMessage(errMsg, lang), true
}

// selectLanguageMessage returns the appropriate message based on the language.
func selectLanguageMessage(msg Message, lang string) string {
	if lang == preference.LANG_EN {
		return msg.EN
	}
	return msg.ID
}

// formatAnnotatedMessage extracts an argument from the error string and formats the message.
func formatAnnotatedMessage(msg, errStr string) string {
	args := fmt.Sprintf("%q", errStr)

	if start, end := strings.LastIndex(args, `{{`), strings.LastIndex(args, `}}`); start > -1 && end > -1 {
		// Extract content between {{ and }}
		args = strings.TrimSpace(args[start+2 : end])
	} else if index := strings.Index(args, `\n`); index > 0 {
		// Extract first line (excluding the opening quote)
		args = strings.TrimSpace(args[1:index])
	}
	// If neither pattern matches, args remains the full quoted string

	return fmt.Sprintf(msg, args)
}
