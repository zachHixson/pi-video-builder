package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

func main(){
	clipSizes, clipPaths := scanClipDir()
	piDigits := getPiDigits()
	generateAllChunks(piDigits, clipSizes, clipPaths)
	fmt.Println("Finished")
}

func getAbsolutePath(path string) string{
	rootDir, rootErr := filepath.Abs(filepath.Dir(os.Args[0]))

	check(rootErr)

	return filepath.Join(rootDir, path)
}

func scanClipDir() ([]int, []string){
	var fileSizes []int
	var filePaths []string

	fmt.Println("Scanning clip directory")

	clipDir := getAbsolutePath(os.Args[1])

	if _, err := os.Stat(clipDir); os.IsNotExist(err) {
		fmt.Println("ERROR: The source clip directory that was provided does not exists. Please check and try again")
		os.Exit(1)
	}

	clipPathErr := filepath.Walk(clipDir, func(path string, info os.FileInfo, err error) error {
		fi, _ := os.Stat(path)
		fileSizes = append(fileSizes, int(fi.Size()))
		filePaths = append(filePaths, path)
		return nil
	})

	check(clipPathErr)

	fileSizes = fileSizes[1:]
	filePaths = filePaths[1:]

	return fileSizes, filePaths
}

func getPiDigits() string{
	fmt.Println("Reading PI digits file")

	filePath := getAbsolutePath(os.Args[2])

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Println("ERROR: The path to the PI text file is invalid. Check path and try again")
		os.Exit(1)
	}

	rawText, err := ioutil.ReadFile(filePath)
	check(err)

	returnString := regexp.MustCompile(`\s|3\.`).ReplaceAllString(string(rawText), "")
	returnString = regexp.MustCompile(`\D`).ReplaceAllString(returnString, "")

	return returnString
}

func generateAllChunks(digits string, sizeArr []int, pathArr []string){
	fmt.Println("Building chunks")

	const CHUNKSIZE = 1300000000
	digitLen := len(digits)
	outDir := getAbsolutePath(os.Args[3]) + "\\"
	i := getResumeDigit()

	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		os.Mkdir(outDir, os.ModeDir)
	}

	if _, err := os.Stat(outDir + "temp.mp4"); !os.IsNotExist(err) {
		fmt.Println("Clearing previous temp file")
		_ = os.Remove(outDir + "temp.mp4")
	}

	for i < digitLen {
		totalBytes := 0
		startDigit := i + 1
		digitChunk := ""
		argClips := ""
		argFilter := ""

		//calculate and display percentage
		perc := float32(i) / float32(len(digits)) * 100
		fmt.Printf("%f%% Current Digit:%d\n", perc, i)

		//build chunk
		for totalBytes < CHUNKSIZE && i < digitLen {
			digitChunk += string(digits[i])
			totalBytes += sizeArr[getCharAsInt(digits, i)]
			i++
		}

		chunkLen := len(digitChunk)
		endDigit := i

		//fill in commands
		for d := 0; d < chunkLen; d++ {
			curDigit := getCharAsInt(digitChunk, d)
			argClips += "-i " + string(pathArr[curDigit]) + " "
			argFilter += "[" + strconv.Itoa(d) + ":v][" + strconv.Itoa(d) + ":a] "
		}

		//execute command
		if len(argClips) > 0 && len(argFilter) > 0 {
			outFileName := strconv.Itoa(startDigit) + "-" + strconv.Itoa(endDigit)
			argStr := "ffmpeg " + argClips
			argStr += "-filter_complex \"" + argFilter
			argStr += "concat=n=" + strconv.Itoa(chunkLen)
			argStr += ":v=1:a=1 [v] [a]\" -map \"[v]\" -map \"[a]\"" + outDir + "temp.mp4"
			pws := exec.Command("powershell", "/c", argStr)
			//pws.Stdout = os.Stdout
			//pws.Stderr = os.Stderr
			err := pws.Run()

			if err != nil {
				fmt.Println("ERROR: Error combining chunk " + outFileName +  " with FFmpeg. Check source clips for corrupted files. If that does not resolve the issue, remove any temp.mp4 files and try again")
				os.Exit(1)
			} else {
				os.Rename(outDir + "temp.mp4", outDir + outFileName + ".mp4")
			}
		}
	}
}

func getCharAsInt(str string, idx int) int{
	c, err := strconv.Atoi(string(str[idx]))

	check(err)

	return c
}

func getResumeDigit() int{
	largestNum := 0
	outPath := getAbsolutePath(os.Args[3])

	if _, err := os.Stat(outPath); os.IsNotExist(err) {
		fmt.Println("No existing clips found, starting at digit 0")
		return 0
	}

	clipPathErr := filepath.Walk(outPath, func(path string, info os.FileInfo, err error) error {
		fileName := path[len(outPath):]
		fileName = regexp.MustCompile(`\\|\.\w+`).ReplaceAllString(fileName, "")
		nums := strings.Split(fileName, "-")

		if len(nums) > 1 {
			bigNum, err := strconv.Atoi(nums[1])

			if err == nil && bigNum > largestNum{
				largestNum = bigNum
			}
		}

		return nil
	})

	check(clipPathErr)

	return largestNum
}

func check(e error){
	if (e != nil) {
		panic(e)
	}
}