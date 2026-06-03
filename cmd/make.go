/*
Copyright © 2026 dayu dayu@dayu.jp
*/

package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dayu255/autocoder/internal/config"
	"github.com/dayu255/autocoder/template"
	"github.com/spf13/cobra"
)

var levels = []string{"a", "b", "c", "d", "e", "f", "g"}

func checkContest(name string) []string {
	if strings.HasPrefix(name, "abc") {
		return levels
	} else if strings.HasPrefix(name, "arc") {
		return levels[:6]
	} else if strings.HasPrefix(name, "agc") {
		return levels[:6]
	} else if strings.HasPrefix(name, "ahc") {
		return levels[:1]
	}

	return levels
}

// 引数の数を確認
func ArgValidation(cmd *cobra.Command, args []string) error {
	switch len(args) {
	case 0:
		return errors.New("contest name is required")
	case 1, 2:
		return nil
	default:
		return errors.New("too many arguments")
	}
}

// makeCmd represents the make command
var makeCmd = &cobra.Command{
	Use:   "make",
	Short: "A brief description of your command",
	Args:  ArgValidation,
	RunE:  Make,
}

func Make(cmd *cobra.Command, args []string) error {
	contest := args[0]
	contestLevels := checkContest(contest)

	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	dirPath := filepath.Join(cwd, contest)

	fileInfo, err := os.Lstat("./")
	if err != nil {
		return fmt.Errorf("failed to stat current directory: %w", err)
	}
	fileMode := fileInfo.Mode()
	unixPerms := fileMode & os.ModePerm

	if err := os.MkdirAll(dirPath, unixPerms); err != nil {
		return fmt.Errorf("failed to create contest directory: %w", err)
	}

	var templateContent string
	var extension string

	cfg, cfgErr := config.LoadConfig()

	language := ""
	if len(args) == 2 {
		language = strings.ToLower(args[1])
	} else if cfg != nil && cfg.DefaultLanguage != "" {
		language = strings.ToLower(cfg.DefaultLanguage)
	} else {
		language = "cpp"
	}

	if cfg != nil {
		for _, tmpl := range cfg.TemplateFile {
			if strings.EqualFold(tmpl.Language, language) {
				data, err := os.ReadFile(tmpl.FilePath)
				if err != nil {
					return fmt.Errorf("failed to read template file: %w", err)
				}
				templateContent = string(data)
				extension = filepath.Ext(tmpl.FilePath)
				break
			}
		}
	}

	if templateContent == "" {
		switch language {
		case "cpp", "c++":
			templateContent = template.DefaultCPP
			extension = ".cpp"
		case "py", "python", "python3":
			templateContent = template.DefaultPY
			extension = ".py"
		}
	}

	if templateContent == "" {
		if cfgErr != nil {
			return fmt.Errorf("no template for language %q (config error: %w)", language, cfgErr)
		}
		return fmt.Errorf("no template for language %q", language)
	}

	if extension == "" {
		extension = "." + language
	}

	var joinedErr error
	for _, level := range contestLevels {
		filePath := filepath.Join(dirPath, level+extension)
		f, err := os.Create(filePath)
		if err != nil {
			joinedErr = errors.Join(joinedErr, fmt.Errorf("failed to create %s: %w", filePath, err))
			continue
		}

		if _, err := f.WriteString(templateContent); err != nil {
			joinedErr = errors.Join(joinedErr, fmt.Errorf("failed to write %s: %w", filePath, err))
		}
		if err := f.Close(); err != nil {
			joinedErr = errors.Join(joinedErr, fmt.Errorf("failed to close %s: %w", filePath, err))
		}
	}

	return joinedErr
}

func init() {
	rootCmd.AddCommand(makeCmd)
}
