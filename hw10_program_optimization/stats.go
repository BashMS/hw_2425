package hw10programoptimization

import (
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type User struct {
	Email string `json:"Email"`
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
	content, err := io.ReadAll(r)
	if err != nil {
		return
	}

	lines := strings.Split(string(content), "\n")
	result = make(users, len(lines))
	for i, line := range lines {
		var user User
		if err = json.Unmarshal([]byte(line), &user); err != nil {
			return
		}
		result[i] = user
	}
	
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	rg := regexp.MustCompile("\\." + domain)
	for _, user := range u {
		matched := rg.MatchString(user.Email)
		if matched {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}
	return result, nil
}
