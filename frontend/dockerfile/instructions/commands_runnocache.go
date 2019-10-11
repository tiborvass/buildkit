// +build dfrunnocache

package instructions

import (
	"github.com/pkg/errors"
)

var noCacheKey = "dockerfile/run/nocache"

func init() {
	parseRunPreHooks = append(parseRunPreHooks, runNoCachePreHook)
	parseRunPostHooks = append(parseRunPostHooks, runNoCachePostHook)
}

func runNoCachePreHook(cmd *RunCommand, req parseRequest) error {
	st := &noCacheState{}
	st.flag = req.flags.AddBool("no-cache", false)
	cmd.setExternalValue(noCacheKey, st)
	return nil
}

func runNoCachePostHook(cmd *RunCommand, req parseRequest) error {
	st := cmd.getExternalValue(noCacheKey).(*noCacheState)
	if st == nil {
		return errors.Errorf("no noCache state")
	}
	st.noCache = st.flag.Value == "true"
	return nil
}

func GetNoCache(cmd *RunCommand) bool {
	return cmd.getExternalValue(noCacheKey).(*noCacheState).noCache
}

type noCacheState struct {
	flag    *Flag
	noCache bool
}
