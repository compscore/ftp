package ftp

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

type optionsStruct struct {
	// Check if file exists
	Exists bool `compscore:"exists"`

	// Check if contents of file matches a substring
	SubstringMatch bool `compscore:"substring_match"`

	// Check if contents of file matches a regex
	RegexMatch bool `compscore:"regex_match"`

	// Check if contents of file matches a string exactly
	Match bool `compscore:"match"`

	// Sha256 hash of the expected output
	Sha256 bool `compscore:"sha256"`

	// Md5 hash of the expected output
	Md5 bool `compscore:"md5"`

	// Sha1 hash of the expected output
	Sha1 bool `compscore:"sha1"`
}

func (e *optionsStruct) Unmarshal(options map[string]interface{}) {
	_, ok := options["exists"]
	if ok {
		e.Exists = true
	}

	_, ok = options["substring_match"]
	if ok {
		e.SubstringMatch = true
	}

	_, ok = options["regex_match"]
	if ok {
		e.RegexMatch = true
	}

	_, ok = options["match"]
	if ok {
		e.Match = true
	}

	_, ok = options["sha256"]
	if ok {
		e.Sha256 = true
	}

	_, ok = options["md5"]
	if ok {
		e.Md5 = true
	}

	_, ok = options["sha1"]
	if ok {
		e.Sha1 = true
	}
}

func (e *optionsStruct) Compare(expectedOutput string, resp *ftp.Response) error {
	if resp == nil {
		return fmt.Errorf("file does not exist")
	}

	bodyBytes, err := io.ReadAll(resp)
	if err != nil {
		return fmt.Errorf("encountered errors while reading body: %s", err)
	}
	body := string(bodyBytes)

	if e.SubstringMatch {
		if !strings.Contains(body, expectedOutput) {
			return fmt.Errorf("substring match mistmatch: execpted \"%s\"", expectedOutput)
		}
	}

	if e.RegexMatch {
		pattern, err := regexp.Compile(expectedOutput)
		if err != nil {
			return fmt.Errorf("invalid regex pattern: \"%s\"; err: \"%s\"", expectedOutput, err)
		}

		if !pattern.MatchString(body) {
			return fmt.Errorf("regex match mismatch: expected \"%s\" got \"%s\"", expectedOutput, body)
		}
	}

	if e.Match {
		if body != expectedOutput {
			return fmt.Errorf("match mismatch: expected \"%s\" got \"%s\"", expectedOutput, body)
		}
	}

	if e.Sha256 {
		hash := fmt.Sprintf("%x", sha256.Sum256(bodyBytes))
		if hash != expectedOutput {
			return fmt.Errorf("sha256 mismatch: expected \"%s\" got \"%s\"", expectedOutput, hash)
		}
	}

	if e.Md5 {
		hash := fmt.Sprintf("%x", md5.Sum(bodyBytes))
		if hash != expectedOutput {
			return fmt.Errorf("md5 mismatch: expected \"%s\" got \"%s\"", expectedOutput, hash)
		}
	}

	if e.Sha1 {
		hash := fmt.Sprintf("%x", sha1.Sum(bodyBytes))
		if hash != expectedOutput {
			return fmt.Errorf("sha1 mismatch: expected \"%s\" got \"%s\"", expectedOutput, hash)
		}
	}

	return nil
}

func Run(ctx context.Context, target string, command string, expectedOutput string, username string, password string, options map[string]interface{}) (bool, string) {
	o := &optionsStruct{}
	o.Unmarshal(options)

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

	err = o.Compare(expectedOutput, resp)
	if err != nil {
		return false, fmt.Sprintf("failed to compare expected output: %s", err)
	}

	return true, ""
}
