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

		if info.IsDir() {
			PrintResult(colors.Gray, convert.Ignore, "IS DIR", info, "")
			return nil
		}

		verdict, words := convert.String(info.Name())
		if verdict == convert.Ignore {
			PrintResult(colors.YellowBright, verdict, "IGNORED", info, "")
			return nil
		}

		newName := strings.Join(words, "_")

		err = os.Rename(filepath.Join(settings.Path, info.Name()), filepath.Join(settings.Path, newName))
		if err != nil {
			return fmt.Errorf("unable to rename '%s' to '%s': %w", info.Name(), newName, err)
		}

		PrintResult(colors.GreenBright, verdict, "RENAMED", info, newName)
		return nil
	})
}

func PrintResult(color colors.Attribute, detectedCase convert.Case, action string, info os.FileInfo, newName string) {
	var postfix string
	if newName != "" {
		postfix = " to '" + newName + "'"
	}

	var detected string
	if detectedCase != convert.Ignore {
		detected = fmt.Sprintf("(%s%s%s)", colors.BlueBright, colors.BlueBright.Reset(), detectedCase.String())
	}

	fmt.Printf("[%s%10s%s] '%s' %s%s\n",
		color, action, color.Reset(),
		info.Name(), detected, postfix)
}
