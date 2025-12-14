package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/gobwas/glob"
)

type FileInfo struct {
	Name  string
	Path  string
	Size  int64
	IsDir bool
}

func shouldIgnore(filePath string, fileName string, ignorePatterns []string) bool {
	for _, pattern := range ignorePatterns {
		g := glob.MustCompile(pattern)
		if g.Match(filePath) || g.Match(fileName) {
			return true
		}

	}
	return false
}

func shouldAccept(filePath string, fileName string, acceptPatterns []string) bool {
	for _, pattern := range acceptPatterns {
		g := glob.MustCompile(pattern)
		if g.Match(filePath) || g.Match(fileName) {
			return true
		}
	}
	return false
}

func getAllEntriesOfDir(dir string, ignorePatterns []string, acceptList []string) ([]FileInfo, error) {
	dirEntry, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var files []FileInfo
	for _, entry := range dirEntry {
		info, err := entry.Info()
		if err != nil {
			return nil, err
		}

		file := FileInfo{
			Name:  entry.Name(),
			Path:  dir + "/" + entry.Name(),
			Size:  info.Size(),
			IsDir: entry.IsDir(),
		}

		checkFilePath := file.Path
		if entry.IsDir() {
			checkFilePath += "/"
		}

		checkFileName := file.Name
		if entry.IsDir() {
			checkFileName += "/"
		}

		if !shouldIgnore(checkFilePath, checkFileName, ignorePatterns) && (shouldAccept(checkFilePath, checkFileName, acceptList) || entry.IsDir()) {
			files = append(files, file)
		}

	}
	return files, nil
}

func getAllFiles(dir string, ignorePatterns []string, acceptList []string) ([]string, error) {
	var allFiles []string
	entries, err := getAllEntriesOfDir(dir, ignorePatterns, acceptList)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir {
			subFiles, err := getAllFiles(entry.Path, ignorePatterns, acceptList)
			if err != nil {
				return nil, err
			}
			allFiles = append(allFiles, subFiles...)
		} else {
			allFiles = append(allFiles, entry.Path)
		}
	}
	return allFiles, nil
}

func saveOutput(files []string, currentDir string, outputDir string, overwrite bool) error {
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return err
	}
	for _, filePath := range files {
		relativePath := strings.TrimPrefix(filePath, currentDir)
		outputPath := path.Join(outputDir, relativePath)

		if !overwrite {
			if _, err := os.Stat(outputPath); err == nil {
				return fmt.Errorf("file already exists: %s", outputPath)
			}
		}

		outputDirPath := path.Dir(outputPath)
		err := os.MkdirAll(outputDirPath, os.ModePerm)
		if err != nil {
			return err
		}
		// copy the file
		input, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}
		err = os.WriteFile(outputPath, input, 0644)
		if err != nil {
			return err
		}

	}

	return nil
}

var defaultAcceptPatterns = []string{
	"**/.env*",
	"**/config.yaml",
	"**/config.json",
	"**/secrets.*",
	"**/firebase-*.json",
	"**/*.pem",
	"**/*.key",
	"**/id_rsa*",
	"**/credentials.json",
}

var defaultIgnorePatterns = []string{
	"**/.git/**",
	"**/node_modules/**",
	"**/vendor/**",
	"**/.idea/**",
	"**/.vscode/**",
	"**/dist/**",
	"**/build/**",
	"**/*.log",
}

func main() {
	targetDir := flag.String("t", ".", "target directory to scan")
	ignorePatterns := flag.String("i", "", "comma-separated glob patterns to ignore")
	acceptPatterns := flag.String("a", "", "comma-separated glob patterns to accept")
	outputDir := flag.String("o", "", "output directory")
	overwrite := flag.Bool("w", false, "overwrite existing files")
	flag.Parse()

	ignoreList := defaultIgnorePatterns
	if *ignorePatterns != "" {
		ignoreList = []string{}
		for _, pattern := range strings.Split(*ignorePatterns, ",") {
			if pattern != "" {
				ignoreList = append(ignoreList, pattern)
			}
		}
	}

	acceptList := defaultAcceptPatterns
	if *acceptPatterns != "" {
		acceptList = []string{}
		for _, pattern := range strings.Split(*acceptPatterns, ",") {
			if pattern != "" {
				acceptList = append(acceptList, pattern)
			}
		}
	}

	files, err := getAllFiles(*targetDir, ignoreList, acceptList)
	if err != nil {
		println("Error:", err.Error())
		return
	}
	err = saveOutput(files, *targetDir, *outputDir, *overwrite)
	if err != nil {
		println("Error:", err.Error())
		return
	}
	fmt.Printf("Found %d files\n", len(files))
}
