/*
Copyright © 2026 dayu dayu@dayu.jp
*/
package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"sync"
	"unicode"

	"github.com/dayu255/autocoder/internal/fetch"
	"github.com/dayu255/autocoder/internal/level"
	"github.com/spf13/cobra"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download AtCoder test cases",

	Run: download,
}

var contestNames = []string{"abc", "arc", "agc", "ahc"}

func checkContestName(name string) (string, string, error) {
	var contestName string
	var contestNum string
	for i := 0; i < len(name); i++ {
		if unicode.IsDigit(rune(name[i])) {
			contestName = name[:i]
			contestNum = name[i:]
			if _, err := strconv.Atoi(contestNum); err != nil {
				return "", "", errors.New("Bad format")
			}
			break
		}
		if i == len(name)-1 {
			return "", "", errors.New("Contest Num not found")
		}
	}

	if !slices.Contains(contestNames, contestName) {
		return "", "", errors.New("Contest Name is unknown")
	}

	return contestName, contestNum, nil
}

const atcoderURL string = "https://atcoder.jp/contests/"

func download(cmd *cobra.Command, args []string) {
	var contest string
	if len(args) == 0 {
		var err error
		cwd, err := os.Getwd()

		if err != nil {
			log.Fatal(err)
			return
		}
		contest = filepath.Base(cwd)
	}

	// 今いるディレクトリ名が"{回数}{コンテスト名}"になっているか確認
	contestName, contestNum, err := checkContestName(contest)
	if err != nil {
		log.Fatal(err)
		return
	}

	// コンテスト名(abc, agcなど)からレベルのスライスを取得
	levels, err := level.CheckContest(contestName)
	if err != nil {
		log.Fatal(err)
		return
	}

	curDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
		return
	}
	testDir := filepath.Join(curDir, ".autocoder", "")
	err = os.MkdirAll(testDir, 0755)
	if err != nil {
		log.Fatal(err)
		return
	}

	baseURL, err := url.JoinPath(
		atcoderURL,
		fmt.Sprintf("%s%s", contestName, contestNum),
		"tasks",
		fmt.Sprintf("%s%s_", contestName, contestNum),
	)
	if err != nil {
		log.Fatal(err)
		return
	}

	for _, level := range levels {
		tests, err := fetch.FetchTest(fmt.Sprint(baseURL, level))
		if err != nil {
			fmt.Println(err)
			continue
		}
		if len(tests) == 0 {
			continue
		}

		err = os.MkdirAll(
			filepath.Join(testDir, level),
			0755,
		)
		if err != nil {
			log.Fatal(err)
			return
		}

		var wgsave sync.WaitGroup
		for _, test := range tests {
			wgsave.Add(1)
			go func(test fetch.Test) {
				err := os.WriteFile(
					filepath.Join(
						testDir,
						level,
						fmt.Sprintf("%d_%s.txt", test.Num, test.Type),
					),
					[]byte(test.Content),
					0644,
				)

				if err != nil {
					log.Println(err)
				}
				wgsave.Done()
			}(test)
		}
		wgsave.Wait()
		fmt.Printf("Save %s %s\n", contest, level)
	}
	fmt.Printf("Finish download %s\n", contest)

}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
