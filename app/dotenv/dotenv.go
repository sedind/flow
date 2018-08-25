// Package dotenv is library which loads variables from vile and loads
// them to environment for current process
// The TL;DR is that you make a .env file that looks something like
//
//    SOME_ENV_VAR=somevalue
//
// and in your code you can call
//
//    dotenv.Load()
//
// and all environment variables declared in the
// .env file will be available through os.Getenv("SOME_ENV_VAR")
package dotenv

import (
	"bufio"
	"errors"
	"io"
	"os"
	"regexp"
	"strings"
)

var defaultFilenames = []string{".env"}

// Load will read env files and load them into ENV for current process
// Calling dotenv.Load function without arguments will load .env file in the current path
func Load(filenames ...string) (err error) {
	filenames = filenamesOrDefault(filenames)
	for _, filename := range filenames {
		err = loadFile(filename)
		if err != nil {
			return
		}
	}
	return
}

func loadFile(filename string) error {
	envMap, err := readFile(filename)
	if err != nil {
		return err
	}

	currentEnv := map[string]bool{}
	rawEnv := os.Environ()
	for _, rawEnvLine := range rawEnv {
		key := strings.Split(rawEnvLine, "=")[0]
		currentEnv[key] = true
	}

	for key, value := range envMap {
		if !currentEnv[key] {
			os.Setenv(key, value)
		}
	}

	return nil
}

func readFile(filename string) (envMap map[string]string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	return parse(file)
}

// parse reads an env file from io.Reader, returning a map of keys and values.
func parse(r io.Reader) (envMap map[string]string, err error) {
	envMap = make(map[string]string)

	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err = scanner.Err(); err != nil {
		return
	}

	for _, fullLine := range lines {
		if !isIgnoredLine(fullLine) {
			var key, value string
			key, value, err = parseLine(fullLine, envMap)

			if err != nil {
				return
			}
			envMap[key] = value
		}
	}
	return
}

func parseLine(line string, envMap map[string]string) (key string, value string, err error) {
	if len(line) == 0 {
		err = errors.New("zero length string")
		return
	}

	// ditch the comments (but keep quoted hashes)
	if strings.Contains(line, "#") {
		segmentsBetweenHashes := strings.Split(line, "#")
		quotesAreOpen := false
		var segmentsToKeep []string
		for _, segment := range segmentsBetweenHashes {
			if strings.Count(segment, "\"") == 1 || strings.Count(segment, "'") == 1 {
				if quotesAreOpen {
					quotesAreOpen = false
					segmentsToKeep = append(segmentsToKeep, segment)
				} else {
					quotesAreOpen = true
				}
			}

			if len(segmentsToKeep) == 0 || quotesAreOpen {
				segmentsToKeep = append(segmentsToKeep, segment)
			}
		}

		line = strings.Join(segmentsToKeep, "#")
	}

	splitString := strings.SplitN(line, "=", 2)
	if len(splitString) != 2 {
		err = errors.New("Can't separate key from value")
		return
	}

	// Parse the key
	key = splitString[0]
	key = strings.Trim(key, " ")

	// Parse the value
	value = parseValue(splitString[1], envMap)
	return
}

func parseValue(value string, envMap map[string]string) string {

	// trim
	value = strings.Trim(value, " ")

	// check if we've got quoted values or possible escapes
	if len(value) > 1 {
		first := string(value[0:1])
		last := string(value[len(value)-1:])
		if first == last && strings.ContainsAny(first, `"'`) {
			// pull the quotes off the edges
			value = value[1 : len(value)-1]
			// handle escapes
			escapeRegex := regexp.MustCompile(`\\.`)
			value = escapeRegex.ReplaceAllStringFunc(value, func(match string) string {
				c := strings.TrimPrefix(match, `\`)
				switch c {
				case "n":
					return "\n"
				case "r":
					return "\r"
				default:
					return c
				}
			})
		}
	}

	// expand variables
	value = os.Expand(value, func(key string) string {
		if val, ok := envMap[key]; ok {
			return val
		}
		if val, ok := os.LookupEnv(key); ok {
			return val
		}
		return ""
	})
	return value
}

func filenamesOrDefault(filenames []string) []string {
	if len(filenames) == 0 {
		return defaultFilenames
	}
	return filenames
}

func isIgnoredLine(line string) bool {
	trimmedLine := strings.Trim(line, " \n\t")
	return len(trimmedLine) == 0 || strings.HasPrefix(trimmedLine, "#")
}
