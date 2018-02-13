package main

import (
	"strings"
	"time"
)

var (
	gcPause             = 10 * time.Second
	maxOrphanRuntimeSec = 20.0
)

func rmOrphanTerminatingDProcs() {
	for {
		poNstart, err := kubectl(true, "get", "po",
			"--selector=dproctype="+string(DProcTerminating), "-o=custom-columns=:metadata.name,:status.startTime", "--field-selector=status.phase=Running",
			"--no-headers")
		if err != nil {
			debug(err.Error())
		}
		if poNstart != "" {
			poname, start := strings.Split(poNstart, "   ")[0], strings.Split(poNstart, "   ")[1]
			debug("found candidate pod " + poname)
			layout := "2006-01-02T15:04:05Z"
			st, err := time.Parse(layout, start)
			if err != nil {
				warn("couldn't parse start time of pod " + poname)
			}
			now := time.Now()
			diff := now.Sub(st)
			if diff.Seconds() > maxOrphanRuntimeSec {
				debug("found orphaned pod " + poname + " with a start time of " + st.String())
				_, err = kubectl(false, "delete", "pod", poname)
				if err != nil {
					warn("couldn't garbage collect orphaned pod " + poname)
				}
			}
		}
		time.Sleep(gcPause)
	}
}