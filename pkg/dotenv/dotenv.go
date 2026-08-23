package dotenv

import (
	"os"
	"strings"
)

func Load(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.SplitSeq(string(file), "\n")
	for line := range lines {
		splited := strings.SplitN(line, "=", 2)
		// Skip the line that misses the equal symbol.
		if len(splited) < 2 {
			continue
		}

		k := splited[0]
		v := splited[1]

		os.Setenv(k, v)
	}

	return nil
}
