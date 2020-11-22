package hooks

import (
	"context"
	"reflect"
)

type Runner interface {
	RunAfter(ctx context.Context, params Params) error
	RunBefore(ctx context.Context, params Params) error
}

// Params contains environment variables for hook commands. All values
// should be in string format.
type Params struct {
	Remote          string `env:"GLMT_REMOTE"`
	Branch          string `env:"GLMT_BRANCH"`
	MergeRequestURL string `env:"GLMT_MR_URL"`
	Project         string `env:"GLMT_PROJECT"`
	Username        string `env:"GLMT_USERNAME"`
}

// Env specifies the environment by the params. Each entry is of the
// form "key=value".
func (h Params) Env() (env []string) {
	rt := reflect.TypeOf(h)
	rv := reflect.ValueOf(h)

	env = make([]string, rt.NumField())
	for i := 0; i < rt.NumField(); i++ {
		env[i] = rt.Field(i).Tag.Get("env") + "=" + rv.Field(i).String()
	}

	return env
}
