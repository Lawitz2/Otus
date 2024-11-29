package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

const (
	singleBytes = iota
	lines
	jsons
)

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	const algorithm = jsons

	switch algorithm {
	case singleBytes:
		return domainStatBytes(r, domain)
	case lines:
		return domainStatLines(r, domain)
	case jsons:
		return domainStatJSON(r, domain)
	}
	return nil, nil
}

// read individual bytes, designed only to find email addresses
// average speed, low memory consumption.
func domainStatBytes(r io.Reader, domain string) (DomainStat, error) {
	var err error
	var s byte
	var at, dot bool
	builder, altBuilder := strings.Builder{}, strings.Builder{}
	dStat := make(DomainStat)
	rBuf := bufio.NewReader(r)

	for err == nil {
		s, err = rBuf.ReadByte()
		// case order is important
		switch {
		case s == '@':
			at = true
		case at && s == '.':
			dot = true
			builder.WriteByte(s)
		case at && dot && s == '"':
			at, dot = false, false
			if altBuilder.String() == domain {
				builder.WriteString(altBuilder.String())
				dStat[strings.ToLower(builder.String())]++
			}
			builder.Reset()
			altBuilder.Reset()
		case at && dot:
			altBuilder.WriteByte(s)
		case at:
			builder.WriteByte(s)
		}
	}

	if err == io.EOF {
		return dStat, nil
	}
	return nil, err
}

// read line by line, designed to find only email addresses
// higher speed, higher memory consumption.
func domainStatLines(r io.Reader, domain string) (DomainStat, error) {
	var err error
	var b string
	var start, end, dotPos int
	dStat := make(DomainStat)
	rBuf := bufio.NewReader(r)

	for err == nil {
		b, err = rBuf.ReadString('\n')
		start = strings.Index(b, "@") + 1
		end = strings.Index(b[start:], `"`) + start
		dotPos = strings.Index(b[start:end], `.`)
		if b[start+dotPos+1:end] == domain {
			dStat[strings.ToLower(b[start:end])]++
		}
	}

	if err == io.EOF {
		return dStat, nil
	}
	return nil, err
}

// parses JSONs, as a result data other than email can also be used
// the slowest option, average memory consumption
// can be improved using 3rd party json parsers, but since it fits into
// requirements i decided to stick with the standard library option.
func domainStatJSON(r io.Reader, domain string) (DomainStat, error) {
	dStat := make(DomainStat)
	decoder := json.NewDecoder(r)
	user := &User{}
	var err error
	for {
		err = decoder.Decode(user)
		if err == io.EOF {
			return dStat, nil
		}
		if err != nil {
			return nil, err
		}
		if strings.Contains(user.Email, "."+domain) {
			dStat[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}
}
