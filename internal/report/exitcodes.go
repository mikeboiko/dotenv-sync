package report

import (
	"errors"
	"fmt"
)

const (
	ExitSuccess     = 0
	ExitOperational = 1
	ExitValidation  = 2
)

type AppError struct {
	Code     string
	ExitCode int
	Problem  string
	Impact   string
	Action   string
	Err      error
}

func (e *AppError) Error() string {
	if e == nil {
		return ""
	}
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Problem, e.Err)
	}
	return e.Problem
}

func (e *AppError) Unwrap() error { return e.Err }

type ExitSignal struct {
	Code int
}

func (e ExitSignal) Error() string {
	return fmt.Sprintf("exit %d", e.Code)
}

func NewAppError(code string, exitCode int, problem, impact, action string, err error) *AppError {
	return &AppError{Code: code, ExitCode: exitCode, Problem: problem, Impact: impact, Action: action, Err: err}
}

func ExitCode(err error) int {
	if err == nil {
		return ExitSuccess
	}
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr.ExitCode
	}
	var signal ExitSignal
	if errors.As(err, &signal) {
		return signal.Code
	}
	return ExitOperational
}

func SilentExit(code int) error {
	if code == ExitSuccess {
		return nil
	}
	return ExitSignal{Code: code}
}
