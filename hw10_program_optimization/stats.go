package hw10programoptimization

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"github.com/valyala/fastjson"
)

type User struct {
	ID       int    `json:"-"`
	Name     string `json:"-"`
	Username string `json:"-"`
	Email    string `json:"Email"`
	Phone    string `json:"-"`
	Password string `json:"-"`
	Address  string `json:"-"`
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
		user.Email = fastjson.GetString([]byte(line), "Email")
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
