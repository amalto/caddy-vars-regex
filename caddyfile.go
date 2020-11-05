package varsregex

import (
	"fmt"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	httpcaddyfile.RegisterHandlerDirective("vars_regex", parseCaddyFileVarsRegex)
}

func parseCaddyFileVarsRegex(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var vrx VarsRegex
	vrx.Overwrite = true
	for h.Next() {
		for nesting := h.Nesting(); h.NextBlock(nesting); {
			rootDirective := h.Val()
			switch rootDirective {
			case "name":
				args := h.RemainingArgs()
				err := singleArgumentCheck(rootDirective, args)
				if err != nil {
					return nil, err
				}
				vrx.Name = args[0]
			case "source":
				args := h.RemainingArgs()
				err := singleArgumentCheck(rootDirective, args)
				if err != nil {
					return nil, err
				}
				vrx.Source = args[0]
			case "pattern":
				args := h.RemainingArgs()
				err := singleArgumentCheck(rootDirective, args)
				if err != nil {
					return nil, err
				}
				vrx.Pattern = args[0]
			case "overwrite":
				args := h.RemainingArgs()
				err := singleArgumentCheck(rootDirective, args)
				if err != nil {
					return nil, err
				}
				if args[0] == "false" {
					vrx.Overwrite = false
				}
			}
		}
	}
	if len(vrx.Name) == 0 || len(vrx.Source) == 0 || len(vrx.Pattern) == 0 {
		return nil, fmt.Errorf("argument count mismatch.  name, soutce and patter must be supplied")
	}
	return vrx, nil
}

func singleArgumentCheck(directive string, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("%s argument has no value", directive)
	}
	if len(args) != 1 {
		return fmt.Errorf("%s argument value of %s is unsupported", directive, args[0])
	}
	return nil
}
