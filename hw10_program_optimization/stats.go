package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/valyala/fastjson" //nolint: depguard
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return countDomains(r, domain)
}

func countDomains(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var user User
		user.Email = fastjson.GetString(scanner.Bytes(), "Email")
		matched := strings.HasSuffix(user.Email, "."+domain)
		if matched {
			spits := strings.SplitN(user.Email, "@", 2)
			if len(spits) != 2 {
				return nil, fmt.Errorf("invalid email: %s", user.Email)
			}
			result[strings.ToLower(spits[1])]++
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
