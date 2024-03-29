package middleware

import (
	"fmt"
	"net/http"
	"runtime"

	echo "github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

type (
	// RecoverConfig defines the config for Recover middleware.
	RecoverConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper echoMiddleware.Skipper

		// Size of the stack to be printed.
		// Optional. Default value 4KB.
		StackSize int `yaml:"stack_size"`

		// DisableStackAll disables formatting stack traces of all other goroutines
		// into buffer after the trace for the current goroutine.
		// Optional. Default value false.
		DisableStackAll bool `yaml:"disable_stack_all"`

		// DisablePrintStack disables printing stack trace.
		// Optional. Default value as false.
		DisablePrintStack bool `yaml:"disable_print_stack"`

		// LogLevel is log level to printing stack trace.
		// Optional. Default value 0 (Print).
		LogLevel log.Lvl
	}
)

var (
	// DefaultRecoverConfig is the default Recover middleware config.
	DefaultRecoverConfig = RecoverConfig{
		Skipper:           echoMiddleware.DefaultSkipper,
		StackSize:         4 << 10, // 4 KB
		DisableStackAll:   false,
		DisablePrintStack: false,
		LogLevel:          0,
	}
)

// Recover returns a middleware which recovers from panics anywhere in the chain
// and handles the control to the centralized HTTPErrorHandler.
func Recover() echo.MiddlewareFunc {
	return RecoverWithConfig(DefaultRecoverConfig)
}

// RecoverWithConfig returns a Recover middleware with config.
// See: `Recover()`.
func RecoverWithConfig(config RecoverConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultRecoverConfig.Skipper
	}
	if config.StackSize == 0 {
		config.StackSize = DefaultRecoverConfig.StackSize
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					stack := make([]byte, config.StackSize)
					length := runtime.Stack(stack, !config.DisableStackAll)
					if !config.DisablePrintStack {
						msg := fmt.Sprintf("[PANIC RECOVER] %v %s\n", err, stack[:length])
						switch config.LogLevel {
						case log.DEBUG:
							c.Logger().Debug(msg)
						case log.INFO:
							c.Logger().Info(msg)
						case log.WARN:
							c.Logger().Warn(msg)
						case log.ERROR:
							c.Logger().Error(msg)
						case log.OFF:
							// None.
						default:
							c.Logger().Print(msg)
						}
					}
					c.Error(err)
				}
			}()
			return next(c)
		}
	}
}

// RecoverWithReturnMsg returns a Recover middleware with http status 200 and custom message.
// See: `Recover()`.
func RecoverWithReturnMsg(msg interface{}) echo.MiddlewareFunc {
	config := DefaultRecoverConfig
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultRecoverConfig.Skipper
	}
	if config.StackSize == 0 {
		config.StackSize = DefaultRecoverConfig.StackSize
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					stack := make([]byte, config.StackSize)
					length := runtime.Stack(stack, !config.DisableStackAll)
					if !config.DisablePrintStack {
						msg := fmt.Sprintf("[PANIC RECOVER] %v %s\n", err, stack[:length])
						switch config.LogLevel {
						case log.DEBUG:
							c.Logger().Debug(msg)
						case log.INFO:
							c.Logger().Info(msg)
						case log.WARN:
							c.Logger().Warn(msg)
						case log.ERROR:
							c.Logger().Error(msg)
						case log.OFF:
							// None.
						default:
							c.Logger().Print(msg)
						}
					}
					_ = c.JSON(http.StatusOK, msg)
				}
			}()
			return next(c)
		}
	}
}

// RecoverWithCustomConfig returns a Recover middleware with http status 200 and custom message.
// See: `Recover()`.
func RecoverWithCustomConfig(config RecoverConfig, msg interface{}) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultRecoverConfig.Skipper
	}
	if config.StackSize == 0 {
		config.StackSize = DefaultRecoverConfig.StackSize
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			defer func() {
				if r := recover(); r != nil {
					err, ok := r.(error)
					if !ok {
						err = fmt.Errorf("%v", r)
					}
					stack := make([]byte, config.StackSize)
					length := runtime.Stack(stack, !config.DisableStackAll)
					if !config.DisablePrintStack {
						msg := fmt.Sprintf("[PANIC RECOVER] %v %s\n", err, stack[:length])
						switch config.LogLevel {
						case log.DEBUG:
							c.Logger().Debug(msg)
						case log.INFO:
							c.Logger().Info(msg)
						case log.WARN:
							c.Logger().Warn(msg)
						case log.ERROR:
							c.Logger().Error(msg)
						case log.OFF:
							// None.
						default:
							c.Logger().Print(msg)
						}
					}
					c.JSON(http.StatusOK, msg)
				}
			}()
			return next(c)
		}
	}
}
