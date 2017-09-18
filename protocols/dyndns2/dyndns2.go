package dyndns2

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/g0dsCookie/ddnsd/log"
)

var regex *regexp.Regexp

func init() {
	regex = regexp.MustCompile(`^(?P<response>(badauth|!donator|good|nochg|notfqdn|nohost|numhost|abuse|badagent|dnserr|911))(?P<args>.*)$`)
}

type Config interface {
	GetUpdateURL() string
	GetHosts() []string
	GetCredentials() (string, string)
	Disable()
	TempDisable()
}

type disable struct{ msg string }

type tempDisable struct{ msg string }

func (e disable) Error() string     { return e.msg }
func (e tempDisable) Error() string { return e.msg }

func update(baseURL string, hosts []string) error {
	joined := strings.Join(hosts, ",")
	log.Debug(fmt.Sprintf("BaseURL: %v", baseURL), fmt.Sprintf("Hosts: %v", joined))
	resp, err := http.Get(fmt.Sprintf("%v%v", baseURL, joined))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(b)

	for i := 0; i < len(hosts); i++ {
		line, err := buf.ReadString('\n')
		log.Debug(fmt.Sprintf("Parsing line '%v'", line))

		n1 := regex.SubexpNames()
		match := regex.FindAllStringSubmatch(line, -1)
		if match == nil {
			log.Err(fmt.Sprintf("could not parse response: %v", line))
			continue
		}
		result := map[string]string{}
		for j, n := range match[0] {
			result[n1[j]] = n
		}

		switch result["response"] {
		case "good":
			if result["args"] == "127.0.0.1" {
				log.Warn(fmt.Sprintf("%v has been updated to 127.0.0.1", hosts[i]))
			} else {
				log.Notice(fmt.Sprintf("%v has been updated to %v", hosts[i], result["args"]))
			}
		case "nochg":
			log.Notice(fmt.Sprintf("%v is already up2date", hosts[i]))

		case "notfqdn":
			log.Crit(fmt.Sprintf("%v has an invalid format", hosts[i]))
			return disable{"invalid hostname"}
		case "nohost":
			log.Crit(fmt.Sprintf("%v is not registered for your account", hosts[i]))
			return disable{"hostname not available"}
		case "numhost":
			log.Crit(fmt.Sprintf("%v is a round robin record which is not allowed", hosts[i]))
			return disable{"round robin record found"}
		case "abuse":
			log.Crit(fmt.Sprintf("%v is blocked for abuse", result["args"]))
			return disable{"blocked for abuse"}

		case "badagent":
			log.Emerg("something has gone terrible wrong")
			os.Exit(1)

		case "dnserr":
			log.Crit("dns error detected")
			return tempDisable{"dnserr returned"}
		case "911":
			log.Crit("problem/maintenance detected")
			return tempDisable{"911 returned"}

		case "badauth":
			log.Crit("invalid username/password specified")
			return disable{"bad authentication"}
		case "!donator":
			log.Crit("tried to use an option only available for credited users")
			return disable{"bad option"}
		}

		if err != nil {
			break
		}
	}

	return nil
}

func Update(c Config) error {
	username, password := c.GetCredentials()
	hosts := c.GetHosts()
	updateURL := c.GetUpdateURL()

	url := bytes.Buffer{}
	url.WriteString("http://")
	url.WriteString(username)
	url.WriteByte(':')
	url.WriteString(password)
	url.WriteByte('@')
	url.WriteString(updateURL)
	url.WriteString("/nic/update?hostname=")

	for len(hosts) > 20 {
		if err := update(url.String(), hosts[:20]); err != nil {
			switch err.(type) {
			case disable:
				c.Disable()
			case tempDisable:
				c.TempDisable()
			}
			return err
		}
		hosts = hosts[20:]
	}

	if len(hosts) > 0 {
		if err := update(url.String(), hosts); err != nil {
			return err
		}
	}

	return nil
}
