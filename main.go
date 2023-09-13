package main

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/jlaffaye/ftp"
)

type expectedOutputStruct struct {
	// Check if file exists
	Exists bool `compscore:"exists"`

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

func (e *expectedOutputStruct) Unmarshal(options map[string]interface{}) error {
	_, ok := options["exists"]
	if ok {
		e.Exists = true
	}

	substringMatchInterface, ok := options["substring_match"]
	if ok {
		substringMatch, ok := substringMatchInterface.(string)
		if !ok {
			return fmt.Errorf("substring_match must be a string")
		}
		e.SubstringMatch = substringMatch
	}

	regexMatchInterface, ok := options["regex_match"]
	if ok {
		regexMatch, ok := regexMatchInterface.(string)
		if !ok {
			return fmt.Errorf("regex_match must be a string")
		}
		e.RegexMatch = regexMatch
	}

	matchInterface, ok := options["match"]
	if ok {
		match, ok := matchInterface.(string)
		if !ok {
			return fmt.Errorf("match must be a string")
		}
		e.Match = match
	}

	sha256Interface, ok := options["sha256"]
	if ok {
		sha256, ok := sha256Interface.(string)
		if !ok {
			return fmt.Errorf("sha256 must be a string")
		}
		e.Sha256 = sha256
	}

	md5Interface, ok := options["md5"]
	if ok {
		md5, ok := md5Interface.(string)
		if !ok {
			return fmt.Errorf("md5 must be a string")
		}
		e.Md5 = md5
	}

	sha1Interface, ok := options["sha1"]
	if ok {
		sha1, ok := sha1Interface.(string)
		if !ok {
			return fmt.Errorf("sha1 must be a string")
		}
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

	if e.Sha256 != "" {
		hash := fmt.Sprintf("%x", sha256.Sum256(bodyBytes))
		if hash != e.Sha256 {
			return fmt.Errorf("sha256 mismatch: expected \"%s\", \"%s\"", e.Sha256, hash)
		}
	}

	if e.Md5 != "" {
		hash := fmt.Sprintf("%x", md5.Sum(bodyBytes))
		if hash != e.Md5 {
			return fmt.Errorf("md5 mismatch: expected \"%s\", \"%s\"", e.Md5, hash)
		}
	}

	if e.Sha1 != "" {
		hash := fmt.Sprintf("%x", sha1.Sum(bodyBytes))
		if hash != e.Sha1 {
			return fmt.Errorf("sha1 mismatch: expected \"%s\", \"%s\"", e.Sha1, hash)
		}
	}

	return nil
}

func Run(ctx context.Context, target string, command string, expectedOutput string, username string, password string, options map[string]interface{}) (bool, string) {
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

	output := &expectedOutputStruct{}
	err = output.Unmarshal(options)
	if err != nil {
		return false, fmt.Sprintf("failed to parse expected output: %s", err)
	}

	err = output.Compare(resp)
	if err != nil {
		return false, fmt.Sprintf("failed to compare expected output: %s", err)
	}

	return true, ""
}
