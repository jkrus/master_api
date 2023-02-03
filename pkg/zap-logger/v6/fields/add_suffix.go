package fields

import (
	"net/http"
	"regexp"
	"runtime"
	"strings"

	"github.com/gosimple/slug"
	"go.uber.org/zap"
)

// ...
const (
	FunctionNameSuffix = "__function_name__"
	CounterSuffix      = "__counter__"

	suffixSeparator  = "_"
	minimalPartCount = 2
)

var (
	globalCounter counter

	anonymousFunctionPattern = regexp.MustCompile(`\bfunc\d+\b`)
	idInPathPattern          = regexp.MustCompile(`/\d+`)
)

// MakeRequestURLSuffix ...
func MakeRequestURLSuffix(request *http.Request) string {
	suffix := request.URL.Path
	suffix = idInPathPattern.ReplaceAllString(suffix, "")

	return suffix
}

// AddSuffix ...
//
// The skippedStackFrames parameter is used to select a function, whose name will act as a suffix.
// Zero means the function, from which AddSuffix() is called.
//
func AddSuffix(field zap.Field, skippedStackFrames int, suffixes ...string) zap.Field {
	// plus one skipped stack frame for the current function
	if suffix := getSuffix(&globalCounter, skippedStackFrames+1, suffixes); len(suffix) != 0 {
		field.Key += suffixSeparator + slug.Make(suffix)
	}

	return field
}

func getSuffix(counter *counter, skippedStackFrames int, suffixes []string) string {
	if len(suffixes) == 0 {
		return ""
	}

	var suffix string
	switch suffixes[0] {
	case FunctionNameSuffix:
		// plus one skipped stack frame for the current function
		name, ok := getFunctionName(skippedStackFrames + 1)
		if ok {
			suffix = name
		} else {
			// plus one skipped stack frame for the recursive call
			suffix = getSuffix(counter, skippedStackFrames+1, suffixes[1:])
		}
	case CounterSuffix:
		suffix = counter.text()
		counter.increment()
	default:
		suffix = suffixes[0]
	}

	return suffix
}

func getFunctionName(skippedStackFrames int) (name string, ok bool) {
	// plus one skipped stack frame for the current function
	pc, _, _, ok := runtime.Caller(skippedStackFrames + 1)
	if !ok {
		return "", false
	}

	function := runtime.FuncForPC(pc)
	if function == nil {
		return "", false
	}

	packageName := getLastParts(function.Name(), "/", 1)
	partCount := minimalPartCount +
		len(anonymousFunctionPattern.FindAllString(packageName, -1))
	return getLastParts(packageName, ".", partCount), true
}

func getLastParts(text string, separator string, n int) string {
	parts := strings.Split(text, separator)
	if partCount := len(parts); partCount > n {
		parts = parts[partCount-n:]
	}

	return strings.Join(parts, separator)
}
