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
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}

	return countDomains(u, domain)
}

type users []User

func getUsers(r io.Reader) (result users, err error) {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		var user User
		user.Email = fastjson.GetString(scanner.Bytes(), "Email")
		result = append(result, user)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	for _, user := range u {
		if !strings.Contains(user.Email, "@") {
			return nil, fmt.Errorf("invalid email: %s", user.Email)
		}
		matched := strings.HasSuffix(user.Email, "."+domain)
		if matched {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}

	return result, nil
}
