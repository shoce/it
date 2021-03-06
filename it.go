/*
it = internet time in beats, there are one thousand beats between midnights

usage: it [options]
options:
	y = print year since last millenium
	d = print day number since last january first
	s = print subbeats (there are one hundred subbeats in one beat)
	n = not to print beats
	f = yd = print year, day and beats
	ff = yds = print year, day, beats and subbeats
	v = print string of year and day suitable for version tagging like 20.289
	vv = print string of year, day and beats suitable for version tagging like 20.289.499

history:
2020/8/31 copy from tik
20/289 (2020/10/15) add options for version printing and not printing beats but date only
20/209 add options f and ff

GoFmt GoBuildNull GoBuild GoRelease GoRun
*/

package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	Beat time.Duration = time.Duration(24) * time.Hour / 1000
)

var (
	tzBiel   *time.Location
	tnow     time.Time
	td0, ty0 time.Time
)

func beatsfmt(td time.Duration) string {
	return fmt.Sprintf("%03d", int(td/Beat))
}

func subbeatsfmt(td time.Duration) string {
	return fmt.Sprintf("%02d", int((td/(Beat/100))%100))
}

func dayfmt(td time.Duration) string {
	day := int(td/(time.Duration(24)*time.Hour)) + 1
	return fmt.Sprintf("%03d", day)
}

func monthdayfmt(td time.Duration) string {
	return ty0.Add(td).Format("0102")
}

func yearfmt(t time.Time, zeropad bool) (y string) {
	if zeropad {
		y = fmt.Sprintf("%03d", t.Year()%1000)
	} else {
		y = fmt.Sprintf("%d", t.Year()%1000)
	}
	return y
}

func init() {
	tzBiel = time.FixedZone("Biel", 60*60)
	tnow = time.Now().In(tzBiel)
	td0 = time.Date(tnow.Year(), tnow.Month(), tnow.Day(), 0, 0, 0, 0, tzBiel)
	ty0 = time.Date(tnow.Year(), 1, 1, 0, 0, 0, 0, tzBiel)
}

func main() {
	var pbeats, psubbeats bool
	var pday, pmonthday, pyear bool
	var pversion, plongversion bool
	pbeats = true
	for _, arg := range os.Args[1:] {
		if strings.Contains(arg, "y") {
			pyear = true
		}
		if strings.Contains(arg, "d") {
			pday = true
			pmonthday = false
		}
		if strings.Contains(arg, "m") {
			pmonthday = true
			pday = false
		}
		if strings.Contains(arg, "s") {
			psubbeats = true
		}
		if strings.Contains(arg, "n") {
			pbeats = false
		}
		if strings.Contains(arg, "f") {
			pyear = true
			if !pmonthday {
				pday = true
			}
			pbeats = true
			psubbeats = false
		}
		if strings.Count(arg, "f") > 1 {
			psubbeats = true
		}
		if strings.Contains(arg, "v") {
			pversion = true
		}
		if strings.Count(arg, "v") > 1 {
			plongversion = true
		}
	}

	td := time.Since(td0)

	var version string
	if pmonthday {
		version = fmt.Sprintf("%s.%s", yearfmt(tnow, true), monthdayfmt(time.Since(ty0)))
	} else {
		version = fmt.Sprintf("%s.%s", yearfmt(tnow, true), dayfmt(time.Since(ty0)))
	}
	if plongversion {
		version += "." + beatsfmt(td)
	}
	if pversion || plongversion {
		fmt.Println(version)
		return
	}

	var s string
	if pbeats {
		s += "@" + beatsfmt(td)
	}
	if psubbeats {
		s += "." + subbeatsfmt(td)
	}
	if pday {
		s = dayfmt(time.Since(ty0)) + s
	} else if pmonthday {
		s = monthdayfmt(time.Since(ty0)) + s
	}
	if pyear {
		s = yearfmt(tnow, false) + "/" + s
	}

	fmt.Println(s)
}
