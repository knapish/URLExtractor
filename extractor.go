package urlextractor

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/idna"
	"golang.org/x/net/publicsuffix"
)

type url struct {
	Scheme     string
	SubDomain  string
	ApexDomain string
	TLD        string
	Port       int
	Path       string
	UserInfo   string
	Query      string
	Fragment   string
}

var apexDomainRegex = `^([a-z0-9]{1})([a-z0-9\-]{0,60})([a-z0-9]{0,1})$`
var subDomainRegex = `^([a-z0-9\_]{1})([a-z0-9\-]{0,60})([a-z0-9]{0,1})$`

func getScheme(s string) (string, string) {
	if strings.Contains(s, "://") {
		schemaParts := strings.Split(s, "://")
		if len(schemaParts) != 2 {
			return s, ""
		}
		return schemaParts[1], schemaParts[0]
	}
	return s, ""
}

func domainPort(s string) (string, int) {
	if strings.Contains(s, ":") {
		portParts := strings.SplitN(s, ":", 2)
		pInt, _ := strconv.Atoi(portParts[1])
		if pInt != 0 {
			if pInt >= 0 && pInt <= 65535 {
				return portParts[0], pInt
			}
		}
	}

	return s, 0
}

func domainPath(s string) (string, string) {
	if strings.Contains(s, "/") {
		pathParts := strings.SplitN(s, "/", 2)
		if pathParts[1] != "" {
			return pathParts[0], pathParts[1]
		}
	}
	return s, ""
}

func userInfo(s string) (string, string) {
	if strings.Contains(s, "@") {
		userInfoParts := strings.Split(s, "@")
		if len(userInfoParts) != 2 {
			lastPart := userInfoParts[len(userInfoParts)-1]
			userInfo := strings.TrimSuffix(s, lastPart)
			return lastPart, userInfo
		} else {
			return userInfoParts[1], userInfoParts[0]
		}
	}

	return s, ""
}

func validateApexDomain(s string) bool {
	re := regexp.MustCompile(apexDomainRegex)
	return re.MatchString(s)
}

func validateSubDomain(s string) bool {
	re := regexp.MustCompile(subDomainRegex)
	if strings.Contains(s, ".") {
		sParts := strings.Split(s, ".")
		for _, part := range sParts {
			if !re.MatchString(part) {
				return false
			}
		}
	} else {
		if !re.MatchString(s) {
			return false
		}
	}

	return true
}

func Extract(s string) (url, error) {
	d := url{}
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ToLower(s)

	// Extract schema from input
	s, d.Scheme = getScheme(s)

	if len(s) < 4 {
		return d, fmt.Errorf("invalid domain - string too short: %s", s)
	}

	// Extract fragment from input
	if strings.Contains(s, "#") {
		fragmentParts := strings.SplitN(s, "#", 2)
		if fragmentParts[1] != "" {
			d.Fragment = fragmentParts[1]
			s = fragmentParts[0]
		}
	}

	// Extract query from input
	if strings.Contains(s, "?") {
		queryParts := strings.SplitN(s, "?", 2)
		if queryParts[1] != "" {
			d.Query = queryParts[1]
			s = queryParts[0]
		}
	}

	// Extract path from input
	s, d.Path = domainPath(s)

	// Extract user info from input
	s, d.UserInfo = userInfo(s)

	// Extract port from input
	s, d.Port = domainPort(s)

	// punycode conversion
	p := idna.New()
	s, err := p.ToASCII(s)
	if err != nil {
		return d, err
	}

	parts := strings.Split(s, ".")

	if len(parts) < 2 {
		return d, fmt.Errorf("invalid domain - missing TLD in string: %s", s)
	}

	// Extract TLD and apex domain from input
	if len(parts) > 3 {
		pTLD := parts[len(parts)-3] + " . " + parts[len(parts)-2] + "." + parts[len(parts)-1]
		tmpTLD, isICANN := publicsuffix.PublicSuffix(pTLD)
		if isICANN && tmpTLD == pTLD {
			if parts[len(parts)-4] != "" {
				d.TLD = pTLD
				d.ApexDomain = parts[len(parts)-4]
			}
		}
	}

	if len(parts) > 2 && d.TLD == "" {
		pTLD := parts[len(parts)-2] + "." + parts[len(parts)-1]
		tmpTLD, isICANN := publicsuffix.PublicSuffix(pTLD)
		if isICANN && tmpTLD == pTLD {
			if parts[len(parts)-3] != "" {
				d.TLD = pTLD
				d.ApexDomain = parts[len(parts)-3]
			}
		}
	}

	if d.TLD == "" {
		pTLD := parts[len(parts)-1]
		tmpTLD, isICANN := publicsuffix.PublicSuffix(pTLD)
		if isICANN && tmpTLD == pTLD {
			if parts[len(parts)-2] != "" {
				d.TLD = pTLD
				d.ApexDomain = parts[len(parts)-2]
			}
		}
	}

	if d.TLD == "" {
		return d, fmt.Errorf("invalid domain - missing valid TLD in string: %s", s)
	}
	// End extract apex domain and TLD

	if !validateApexDomain(d.ApexDomain) {
		return d, fmt.Errorf("invalid domain - invalid apex domain: %s", d.ApexDomain)
	}

	d.SubDomain = strings.TrimSuffix(s, "."+d.ApexDomain+"."+d.TLD)
	if d.SubDomain == s {
		d.SubDomain = ""
	}

	if d.SubDomain != "" && !validateSubDomain(d.SubDomain) {
		return d, fmt.Errorf("invalid domain - invalid sub domain: %s", d.SubDomain)
	}

	return d, nil
}
