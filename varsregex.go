package varsregex

import (
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

func init() {
	caddy.RegisterModule(VarsRegex{})
}

type VarsRegex struct {
	Name      string `json:"name,omitempty"`
	Source    string `json:"source,omitempty"`
	Pattern   string `json:"pattern,omitempty"`
	Overwrite bool   `json:"overwrite,omitempty"`

	logger *zap.Logger

	compiled *regexp.Regexp

	repl *caddy.Replacer
}

// Variables used for replacing Caddy placeholders in Source
var (
	placeholderRegexp = regexp.MustCompile(`{([\w.-]+)}`)
)

const rootPlaceholder = "http.vars_regex."

func (VarsRegex) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.vars_regex",
		New: func() caddy.Module { return new(VarsRegex) },
	}
}

func (vrx *VarsRegex) Provision(ctx caddy.Context) error {
	var err error
	vrx.logger = ctx.Logger(vrx)
	vrx.compiled, err = regexp.Compile(vrx.Pattern)
	if err != nil {
		return fmt.Errorf("invalid pattern detected %v", err)
	}
	return nil
}

func (vrx VarsRegex) placeHolderExpansion(wrappedPlaceholderName string) string {
	placeholderName := strings.TrimSuffix(strings.TrimPrefix(wrappedPlaceholderName, "{"), "}")
	ph, found := vrx.repl.GetString(placeholderName)
	if found {
		return ph
	}
	return wrappedPlaceholderName
}

func (vrx VarsRegex) assignPlaceHolder(name string, value string) {
	currentValue, hasValue := vrx.repl.GetString(name)
	if (hasValue && len(currentValue) > 0) && !vrx.Overwrite {
		vrx.logger.Debug("Skipping placeholder assignment, value already set: ", zap.String(name, currentValue))
	} else {
		vrx.repl.Set(name, value)
		vrx.logger.Debug("Adding placeholder: ", zap.String(name, value))
	}
}

func (vrx VarsRegex) ServeHTTP(resp http.ResponseWriter, req *http.Request, next caddyhttp.Handler) error {

	repl := req.Context().Value(caddy.ReplacerCtxKey).(*caddy.Replacer)
	vrx.repl = repl
	expandedSource := placeholderRegexp.ReplaceAllStringFunc(vrx.Source, vrx.placeHolderExpansion)

	matchGroups := vrx.compiled.FindStringSubmatch(expandedSource)
	for i, name := range vrx.compiled.SubexpNames() {
		if i > 0 && i <= len(matchGroups) {
			placeholder := rootPlaceholder + vrx.Name
			if len(name) > 0 {
				placeholder = placeholder + "." + strings.ToLower(name)
			} else {
				placeholder = placeholder + ".capture_group" + strconv.Itoa(i)
			}
			vrx.assignPlaceHolder(placeholder, matchGroups[i])
		}
	}

	matches := vrx.compiled.FindAllStringSubmatch(expandedSource, -1)
	for i, val := range matches {
		for _, m := range val {
			placeholder := rootPlaceholder + vrx.Name + ".match" + strconv.Itoa(i+1)
			vrx.assignPlaceHolder(placeholder, m)
			if len(matchGroups) > 0 {
				break
			}
		}
	}

	return next.ServeHTTP(resp, req)
}

// Interface guard
var _ caddyhttp.MiddlewareHandler = (*VarsRegex)(nil)
