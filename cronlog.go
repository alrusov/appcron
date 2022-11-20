package appcron

import (
	"bytes"
	"fmt"
	"time"

	"github.com/alrusov/log"
	"github.com/alrusov/misc"
)

//----------------------------------------------------------------------------------------------------------------------------//

// CronLog --
type CronLog struct{}

//----------------------------------------------------------------------------------------------------------------------------//

// Info --
func (cl *CronLog) Info(msg string, keysAndValues ...any) {
	log.Message(log.TRACE2, cl.makeMsg(nil, msg, keysAndValues...))
}

// Error --
func (cl *CronLog) Error(err error, msg string, keysAndValues ...any) {
	log.Message(log.ERR, cl.makeMsg(err, msg, keysAndValues...))
}

//----------------------------------------------------------------------------------------------------------------------------//

func (cl *CronLog) makeMsg(err error, msg string, keysAndValues ...any) string {
	out := new(bytes.Buffer)

	if err != nil {
		out.WriteString(err.Error())
	}

	if msg != "" {
		if out.Len() != 0 {
			out.WriteString(": ")
		}
		out.WriteString(msg)
	}

	kvFmt := cl.makeFmt(keysAndValues...)
	if kvFmt != "" {
		if out.Len() == 0 {
			out.WriteString(kvFmt)
		} else {
			out.WriteString(" (")
			out.WriteString(kvFmt)
			out.WriteString(")")
		}
	}

	return fmt.Sprintf("[cron] "+out.String(), keysAndValues...)
}

//----------------------------------------------------------------------------------------------------------------------------//

func (cl *CronLog) makeFmt(p ...any) string {
	n := len(p)

	if n == 0 {
		return ""
	}

	var fmt = new(bytes.Buffer)

	for i := 0; i < n; i += 2 {
		if i > 0 {
			fmt.WriteString(", ")
		}
		fmt.WriteString("%v=%v")

		i2 := i + 1
		switch p[i2].(type) {
		case time.Time:
			p[i2] = p[i2].(time.Time).Format(misc.DateTimeFormatJSON)
		}
	}

	return fmt.String()
}

//----------------------------------------------------------------------------------------------------------------------------//
