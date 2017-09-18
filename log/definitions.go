package log

import (
	"errors"
	"log/syslog"
	"regexp"
)

type LogTarget string

type LogSeverity string

const (
	LogSeverityEmerg  LogSeverity = "emergency"
	LogSeverityAlert              = "alert"
	LogSeverityCrit               = "critical"
	LogSeverityErr                = "error"
	LogSeverityWarn               = "warning"
	LogSeverityNotice             = "notice"
	LogSeverityInfo               = "info"
	LogSeverityDebug              = "debug"
)

var logSeverityMap = map[LogSeverity]syslog.Priority{
	LogSeverityEmerg:  syslog.LOG_ERR,
	LogSeverityAlert:  syslog.LOG_ALERT,
	LogSeverityCrit:   syslog.LOG_CRIT,
	LogSeverityErr:    syslog.LOG_ERR,
	LogSeverityWarn:   syslog.LOG_WARNING,
	LogSeverityNotice: syslog.LOG_NOTICE,
	LogSeverityInfo:   syslog.LOG_INFO,
	LogSeverityDebug:  syslog.LOG_DEBUG,
}

type LogFacility string

const (
	LogFacilityKern     LogFacility = "kern"
	LogFacilityUser                 = "user"
	LogFacilityMail                 = "mail"
	LogFacilityDaemon               = "daemon"
	LogFacilityAuth                 = "auth"
	LogFacilitySyslog               = "syslog"
	LogFacilityLPR                  = "lpr"
	LogFacilityNews                 = "news"
	LogFacilityUUCP                 = "uucp"
	LogFacilityCron                 = "cron"
	LogFacilityAuthPriv             = "authpriv"
	LogFacilityFTP                  = "ftp"
	LogFacilityLocal0               = "local0"
	LogFacilityLocal1               = "local1"
	LogFacilityLocal2               = "local2"
	LogFacilityLocal3               = "local3"
	LogFacilityLocal4               = "local4"
	LogFacilityLocal5               = "local5"
	LogFacilityLocal6               = "local6"
	LogFacilityLocal7               = "local7"
)

var logFacilityMap = map[LogFacility]syslog.Priority{
	LogFacilityKern:     syslog.LOG_KERN,
	LogFacilityUser:     syslog.LOG_USER,
	LogFacilityMail:     syslog.LOG_MAIL,
	LogFacilityDaemon:   syslog.LOG_DAEMON,
	LogFacilityAuth:     syslog.LOG_AUTH,
	LogTargetSyslog:     syslog.LOG_SYSLOG,
	LogFacilityLPR:      syslog.LOG_LPR,
	LogFacilityNews:     syslog.LOG_NEWS,
	LogFacilityCron:     syslog.LOG_CRON,
	LogFacilityAuthPriv: syslog.LOG_AUTHPRIV,
	LogFacilityFTP:      syslog.LOG_FTP,
	LogFacilityLocal0:   syslog.LOG_LOCAL0,
	LogFacilityLocal1:   syslog.LOG_LOCAL1,
	LogFacilityLocal2:   syslog.LOG_LOCAL2,
	LogFacilityLocal3:   syslog.LOG_LOCAL3,
	LogFacilityLocal4:   syslog.LOG_LOCAL4,
	LogFacilityLocal5:   syslog.LOG_LOCAL5,
	LogFacilityLocal6:   syslog.LOG_LOCAL6,
	LogFacilityLocal7:   syslog.LOG_LOCAL7,
}

func (l LogSeverity) Priority() syslog.Priority {
	if v, ok := logSeverityMap[l]; ok {
		return v
	}
	panic("unknown severity " + string(l))
}

func (l LogSeverity) Lower(r LogSeverity) bool {
	left, right := l.Priority(), r.Priority()
	return left < right
}

func (l LogSeverity) LowerOrEqual(r LogSeverity) bool {
	left, right := l.Priority(), r.Priority()
	return left <= right
}

func (l LogSeverity) Equal(r LogSeverity) bool {
	left, right := l.Priority(), r.Priority()
	return left == right
}

func (l LogSeverity) BiggerOrEqual(r LogSeverity) bool {
	left, right := l.Priority(), r.Priority()
	return left >= right
}

func (l LogSeverity) Bigger(r LogSeverity) bool {
	left, right := l.Priority(), r.Priority()
	return left > right
}

func (l LogFacility) Priority() syslog.Priority {
	if v, ok := logFacilityMap[l]; ok {
		return v
	}
	panic("unknown facility " + string(l))
}

type LogConfig struct {
	Target        LogTarget   `xml:"target"`
	File          string      `xml:"file"`
	SyslogAddress string      `xml:"address"`
	SyslogTag     string      `xml:"tag"`
	Severity      LogSeverity `xml:"severity"`
	Facility      LogFacility `xml:"facility"`
}

func (c LogConfig) Apply() error {
	return Apply(c)
}

func (c LogConfig) ParseSyslogAddress() (string, string, error) {
	const regex = `^(?P<network>(unix|unixgram|unixpacket|tcp|tcp4|tcp6|udp|udp4|udp6))://(?P<address>(((?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)|\[(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))\]):([0-9]{1,4}|[1-5][0-9]{4}|6[0-4][0-9]{3}|65[0-4][0-9]{2}|655[0-2][0-9]|6553[0-5])|[\w\/ ?-]+))$`
	r, err := regexp.Compile(regex)
	if err != nil {
		return "", "", err
	}

	n1 := r.SubexpNames()
	match := r.FindAllStringSubmatch(c.SyslogAddress, -1)
	if match == nil {
		return "", "", errors.New("could not parse syslog address: " + c.SyslogAddress)
	}
	result := map[string]string{}
	for i, n := range match[0] {
		result[n1[i]] = n
	}

	return result["network"], result["address"], nil
}
