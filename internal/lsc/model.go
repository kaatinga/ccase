package lsc

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/kaatinga/ccase/internal/convert"
	"github.com/kaatinga/ccase/internal/settings"

	"github.com/SuperPaintman/nice/colors"
	"github.com/urfave/cli/v2"
)

func UpdateFiles(_ *cli.Context) error {
	return filepath.Walk(settings.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		input := []rune(info.Name())

		if info.IsDir() && input[0] != '.' {
			// PrintResult(colors.Gray, convert.Ignore, "DIRECTORY", info, "")
			return nil
		}

		verdict, words := convert.String(input)
		switch verdict {
		case convert.Ignore:
			// PrintResult(colors.Gray, verdict, "IGNORED", info, "")
			return nil
		case convert.IsNotDotGo:
			// PrintResult(colors.Gray, verdict, "NOT .GO", info, "")
			return nil
		}

		newName := strings.Join(words, "_") + ".go"

		if info.Name() != newName {
			if err := os.Rename(path, filepath.Join(filepath.Dir(path), newName)); err != nil {
				return fmt.Errorf("unable to rename '%s' to '%s': %w", info.Name(), newName, err)
			}

			PrintResult(colors.YellowBright, verdict, "RENAMED", info, newName)
		}

		return nil
	})
}

func PrintResult(color colors.Attribute, detectedCase convert.Case, action string, info os.FileInfo, newName string) {
	var postfix string
	if newName != "" {
		postfix = " to '" + newName + "'"
	}

	var detected string
	if !(detectedCase == convert.Ignore || detectedCase == convert.IsNotDotGo) {
		detected = fmt.Sprintf("(%s%s%s)", colors.Blue, detectedCase.String(), colors.Blue.Reset())
	}

	fmt.Printf("[%s%10s%s] '%s' %s%s\n",
		color, action, color.Reset(),
		info.Name(), detected, postfix)
}
