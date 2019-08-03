package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const GIT_STATUS_CMD = "git status -s | awk '{print $2}'"
const CONFLICT_SIGNAL string = "<<<<<<"

// Log for record conflict info
type Log struct {
	Filename string
	Row      int
	Col      int
	Conflict []string
}

func main() {
	fmt.Println(White("pre-commit: "), Blue("checking..."))

	files := GetGitStatusFiles()
	var logs []Log
	for _, file := range files {
		logs = append(logs, GetLog(file)...)
	}

	if len(logs) == 0 {
		fmt.Println(White("pre-commit: "), Blue("check success :) \n"))
	} else {
		fmt.Println(White("pre-commit: "), Red("check fail \n"))
		PrintLog(logs)
	}

}

// GetGitStatusFiles from git status comand
func GetGitStatusFiles() []string {
	cmd := exec.Command("sh", "-c", GIT_STATUS_CMD)
	out, _ := cmd.Output()
	list := strings.Split(string(out), "\n")
	if len(list) != 0 {
		list = list[:len(list)-1]
	}
	var files []string
	for _, filepath := range list {
		files = append(files, GetFiles(filepath)...)
	}
	return files
}

// GetFiles from path
func GetFiles(filepath string) []string {
	var files []string
	PthSep := string(os.PathSeparator)

	file, err := os.Stat(filepath)
	if err == nil {
		if file.IsDir() {
			dir, _ := ioutil.ReadDir(filepath)
			for _, childFile := range dir {
				childFilePath := filepath + PthSep + childFile.Name()
				childFile, _ := os.Stat(childFilePath)
				if childFile.IsDir() {
					files = append(files, GetFiles(childFilePath)...)
				} else {
					files = append(files, childFilePath)
				}
			}
		} else {
			files = append(files, filepath)
		}
	}
	return files
}

// GetLog conflict info from file
func GetLog(filePath string) []Log {
	var logs []Log
	data, err := ioutil.ReadFile(filePath)
	if err == nil {
		content := string(data)
		lines := strings.Split(content, "\n")
		for i := 0; i < len(lines); i++ {
			if strings.Contains(lines[i], CONFLICT_SIGNAL) {
				var conflict []string
				if i-1 >= 0 {
					conflict = append(conflict, lines[i-1])
				}
				conflict = append(conflict, Red(lines[i]))
				if i+1 < len(lines) {
					conflict = append(conflict, lines[i+1])
				}

				logs = append(logs, Log{
					Filename: filePath,
					Row:      i + 1,
					Col:      strings.Index(lines[i], CONFLICT_SIGNAL) + 1,
					Conflict: conflict,
				})
			}
		}
	}

	return logs
}

// PrintLog log list
func PrintLog(logs []Log) {
	for i := 0; i < len(logs); i++ {
		fmt.Println(RedFlash("Error:"))
		fmt.Println(Grey("      File: ") + Blue(logs[i].Filename))
		fmt.Println(Grey("      Line: ") + Blue(strconv.Itoa(logs[i].Row)))
		fmt.Println(Grey("       Col: ") + Blue(strconv.Itoa(logs[i].Col)))
		fmt.Println(Grey("  Conflict: "))

		lines := logs[i].Conflict
		for _, line := range lines {
			fmt.Println("            " + line)
		}
	}
}

func colorText(color, text string) string {
	return "\u001B[" + color + "m" + text + "\u001B[0m"
}

// Red color
func Red(text string) string {
	return colorText("1;31", text)
}

// Blue color
func Blue(text string) string {
	return colorText("1;34", text)
}

// Grey color
func Grey(text string) string {
	return colorText("2;37", text)
}

// White color
func White(text string) string {
	return colorText("1;37", text)
}

// Magenta color
func Magenta(text string) string {
	return colorText("1;35", text)
}

// RedFlash color
func RedFlash(text string) string {
	return colorText("1;31;5", text)
}
