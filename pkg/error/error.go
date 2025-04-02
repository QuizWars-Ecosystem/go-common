package error

import (
	"errors"
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/google/uuid"
)

type GrpcStatus interface {
	GRPCStatus() *status.Status
}

var (
	_ error      = (*Error)(nil)
	_ GrpcStatus = (*Error)(nil)
)

type Error struct {
	Code       codes.Code
	StackTrace string
	IncidentID string

	innerError error
	hiderError bool
	message    string
}

func (e *Error) GRPCStatus() *status.Status {
	return status.New(e.Code, e.message)
}

func (e *Error) ToGRPCError() error {
	return status.New(e.Code, e.message).Err()
}

func (e *Error) Error() string {
	return e.error(false)
}

func (e *Error) SafeError() string {
	return e.error(true)
}

func (e *Error) GetCode() codes.Code {
	return e.Code
}

func (e *Error) Unwrap() error {
	return e.innerError
}

func (e *Error) error(safe bool) string {
	switch {
	case e.innerError == nil:
		return e.message
	case safe && e.hiderError:
		return e.message
	case e.message == "":
		return e.innerError.Error()
	default:
		return fmt.Sprintf("%s: %s", e.message, e.innerError.Error())
	}
}

func WrapInternal(err error) *Error {
	if err == nil {
		return nil
	}

	return Internal(err)
}

func Internal(err error) *Error {
	appErr := InternalWithoutStackTrace(err)
	appErr.StackTrace = string(debug.Stack())
	return appErr
}

func InternalWithoutStackTrace(err error) *Error {
	appErr := newHiddenError(err, codes.Internal, "internal error")
	appErr.IncidentID = uuid.New().String()
	return appErr
}

func EnsureInternal(err error) error {
	var appErr *Error
	if !errors.As(err, &appErr) {
		return Internal(err)
	}
	return err
}

func BadRequest(err error) *Error {
	return newWrappedError(err, codes.InvalidArgument)
}

func BadRequestHidden(err error, message string) *Error {
	return newHiddenError(err, codes.InvalidArgument, message)
}

func ValidationError(err error) *Error {
	return newWrappedValidationError(err, codes.InvalidArgument)
}

func NotFound(subject, key string, value any) *Error {
	return newError(codes.NotFound, fmt.Sprintf("%s %s: %v not found", subject, key, value))
}

func AlreadyExists(subject, key string, value any) *Error {
	return newError(codes.AlreadyExists, fmt.Sprintf("%s %s: %v already exists", subject, key, value))
}

func Unauthorized(message string) *Error {
	return newError(codes.Unauthenticated, message)
}

func UnauthorizedHidden(err error, message string) *Error {
	return newHiddenError(err, codes.Unauthenticated, message)
}

func Forbidden(message string) *Error {
	return newError(codes.PermissionDenied, message)
}

func VersionMismatch(subject, key string, value any, version int) *Error {
	return newError(codes.InvalidArgument, fmt.Sprintf("stale version %d for %s %s: %v", version, subject, key, value))
}

func newError(code codes.Code, message string) *Error {
	return &Error{
		Code:    code,
		message: message,
	}
}

func newWrappedError(err error, code codes.Code) *Error {
	return &Error{
		Code:       code,
		innerError: err,
	}
}

func newWrappedValidationError(err error, code codes.Code) *Error {
	var vErr *protovalidate.ValidationError
	if errors.As(err, &vErr) {
		builder := &strings.Builder{}
		builder.WriteString("validation error: ")
		if len(vErr.Violations) > 0 {
			if len(vErr.Violations[0].Proto.GetField().GetElements()) > 0 {
				builder.WriteString(vErr.Violations[0].Proto.GetField().GetElements()[0].GetFieldName())
			}
			builder.WriteString(": ")
			builder.WriteString(vErr.Violations[0].Proto.GetMessage())
		}

		return &Error{
			Code:       code,
			innerError: err,
			message:    builder.String(),
		}
	}

	return &Error{
		Code:       code,
		innerError: err,
	}
}

func newHiddenError(err error, code codes.Code, message string) *Error {
	return &Error{
		Code:       code,
		message:    message,
		innerError: err,
		hiderError: true,
	}
}
