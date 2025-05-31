/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// helpCmd represents the help command
var helpCmd = &cobra.Command{
	Use:   "help",
	Short: "도움말 보기",
	Long:  "모든 명령어와 사용법을 요약해서 출력합니다.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(`morama 0.1.0
영화·드라마 기록 CLI

Usage:
  morama [command]

Available Commands:
  add         기록 추가
  delete      기록 삭제
  edit        기록 수정
  list        기록 목록보기
  show        상세 후기 보기
  report      연도/유형별 통계 출력
  help        도움말

Flags:
  -h, --help     help for morama
      --version  버전 정보 출력`)
	},
}

func init() {
	rootCmd.AddCommand(helpCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// helpCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// helpCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
