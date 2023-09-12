package main

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/jlaffaye/ftp"
)

type expectedOutputStruct struct {
	// Check if contents of file matches a substring
	SubstringMatch string `compscore:"substring_match"`

	// Check if contents of file matches a regex
	RegexMatch string `compscore:"regex_match"`

	// Check if contents of file matches a string exactly
	Match string `compscore:"match"`

	// Sha256 hash of the expected output
	Sha256 string `compscore:"sha256"`

	// Md5 hash of the expected output
	Md5 string `compscore:"md5"`

	// Sha1 hash of the expected output
	Sha1 string `compscore:"sha1"`
}

func (e *expectedOutputStruct) Unmarshal(in string) error {
	structLookup := make(map[string]string)
	split := strings.Split(in, ";")
	for _, item := range split {
		itemSplit := strings.Split(item, "=")
		if len(itemSplit) != 2 {
			return fmt.Errorf("invalid parameter string: %s", item)
		}
		structLookup[strings.TrimSpace(itemSplit[0])] = strings.TrimSpace(itemSplit[1])
	}

	substringMatch, ok := structLookup["substring_match"]
	if ok {
		e.SubstringMatch = substringMatch
	}

	regexMatch, ok := structLookup["regex_match"]
	if ok {
		e.RegexMatch = regexMatch
	}

	match, ok := structLookup["match"]
	if ok {
		e.Match = match
	}

	sha256, ok := structLookup["sha256"]
	if ok {
		e.Sha256 = sha256
	}

	md5, ok := structLookup["md5"]
	if ok {
		e.Md5 = md5
	}

	sha1, ok := structLookup["sha1"]
	if ok {
		e.Sha1 = sha1
	}

	return nil
}

func (e *expectedOutputStruct) Compare(resp *ftp.Response) error {
	if resp == nil {
		return fmt.Errorf("file does not exist")
	}

	bodyBytes, err := io.ReadAll(resp)
	if err != nil {
		return fmt.Errorf("encountered errors while reading body: %s", err)
	}
	body := string(bodyBytes)

	if e.SubstringMatch != "" {
		if !strings.Contains(body, e.SubstringMatch) {
			return fmt.Errorf("substring match mistmatch: execpted \"%s\"", e.SubstringMatch)
		}
	}

	if e.RegexMatch != "" {
		pattern, err := regexp.Compile(e.RegexMatch)
		if err != nil {
			return fmt.Errorf("invalid regex pattern: \"%s\"", err)
		}

		if !pattern.MatchString(body) {
			return fmt.Errorf("regex match mitmatch: expected \"%s\"", e.RegexMatch)
		}
	}

	if e.Match != "" {
		if body != e.Match {
			return fmt.Errorf("match mismatch: expected \"%s\", \"%s\"", e.Match, body)
		}
	}

	return nil
}

func Run(ctx context.Context, target string, command string, expectedOutput string, username string, password string) (bool, string) {
	if !strings.Contains(target, ":") {
		target = target + ":21"
	}

	conn, err := ftp.Dial(target, ftp.DialWithContext(ctx))
	if err != nil {
		return false, fmt.Sprintf("failed to connect to target: %s", err)
	}

	if username != "" && password != "" {
		err = conn.Login(username, password)
		if err != nil {
			return false, fmt.Sprintf("failed to login to target: %s", err)
		}
	} else if username != "" {
		err = conn.Login(username, "")
		if err != nil {
			return false, fmt.Sprintf("failed to login to target: %s", err)
		}
	} else {
		err = conn.Login("anonymous", "anonymous")
		if err != nil {
			return false, fmt.Sprintf("failed to login to target: %s", err)
		}
	}

	resp, err := conn.Retr(command)
	if err != nil {
		return false, fmt.Sprintf("failed to retrieve file: %s", err)
	}
	defer resp.Close()

	expectedOutputStruct := &expectedOutputStruct{}
	err = expectedOutputStruct.Unmarshal(expectedOutput)
	if err != nil {
		return false, fmt.Sprintf("failed to parse expected output: %s", err)
	}

	err = expectedOutputStruct.Compare(resp)
	if err != nil {
		return false, fmt.Sprintf("failed to compare expected output: %s", err)
	}

	return true, ""
}
