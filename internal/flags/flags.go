package flags

import (
	"os"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

func MapWithPrefix(f *pflag.FlagSet, name string, errorHandling pflag.ErrorHandling, prefix string) *pflag.FlagSet {
	fNew := pflag.NewFlagSet(name, errorHandling)
	f.VisitAll(func(flag *pflag.Flag) {
		if prefix != "" {
			flag.Name = prefix + "_" + flag.Name
		}
		fNew.AddFlag(flag)
	})

	return fNew
}

func MustBindEnvToFlagSet(fs *pflag.FlagSet) {
	if err := BindEnvToFlagSet(fs); err != nil {
		panic(err)
	}
}

// BindEnvToFlagSet maps env variables values to config flags.
func BindEnvToFlagSet(fs *pflag.FlagSet) error {
	set := map[string]bool{}
	fs.Visit(func(f *pflag.Flag) {
		set[f.Name] = true
	})

	var flagError error
	fs.VisitAll(func(f *pflag.Flag) {
		if flagError != nil {
			return
		}

		replacer := strings.NewReplacer("-", "_", ".", "_")
		envVar := replacer.Replace(strings.ToUpper(f.Name))

		valEnv, ok := os.LookupEnv(envVar)
		if !ok {
			return
		}

		if set[f.Name] {
			return
		}

		t := f.Value.Type()
		if t == "stringArray" || t == "stringSlice" {
			vals := strings.Split(valEnv, " ")
			for _, v := range vals {
				if err := fs.Set(f.Name, v); err != nil {
					flagError = errors.Wrapf(err, "wrapping %q with %v", f.Name, v)
					return
				}
			}

			return
		}

		if err := fs.Set(f.Name, valEnv); err != nil {
			flagError = errors.Wrapf(err, "wrapping %q with %v", f.Name, valEnv)
			return
		}
	})
	return flagError
}
