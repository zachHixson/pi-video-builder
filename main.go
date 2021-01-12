package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"io/ioutil"
	"regexp"
	"strconv"
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

	rawText, err := ioutil.ReadFile(getAbsolutePath(os.Args[2]))
	check(err)

	returnString := regexp.MustCompile(`\s|3\.`).ReplaceAllString(string(rawText), "")
	returnString = regexp.MustCompile(`\D`).ReplaceAllString(returnString, "")
	fmt.Println(returnString)

	return returnString
}

func generateAllChunks(digits string, sizeArr []int, pathArr []string){
	fmt.Println("Building chunks")

	const CHUNKSIZE = 1300000000
	digitLen := len(digits)
	chunkNum, i := getLogData()
	outDir := getAbsolutePath(os.Args[3]) + "\\"

	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		os.Mkdir(outDir, os.ModeDir)
	}

	for i < digitLen {
		totalBytes := 0
		digitChunk := ""
		argClips := ""
		argFilter := ""

		writeLogFile(chunkNum, i)

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

		//fill in commands
		for d := 0; d < chunkLen; d++ {
			curDigit := getCharAsInt(digitChunk, d)
			argClips += "-i " + string(pathArr[curDigit]) + " "
			argFilter += "[" + strconv.Itoa(d) + ":v][" + strconv.Itoa(d) + ":a] "
		}

		//execute command
		if len(argClips) > 0 && len(argFilter) > 0 {
			argStr := "ffmpeg " + argClips
			argStr += "-filter_complex \"" + argFilter
			argStr += "concat=n=" + strconv.Itoa(chunkLen)
			argStr += ":v=1:a=1 [v] [a]\" -map \"[v]\" -map \"[a]\"" + outDir + strconv.Itoa(chunkNum) + ".mp4"
			pws := exec.Command("powershell", "/c", argStr)
			//pws.Stdout = os.Stdout
			//pws.Stderr = os.Stderr
			err := pws.Run()

			check(err)

			chunkNum++
		}
	}
}

func getCharAsInt(str string, idx int) int{
	c, err := strconv.Atoi(string(str[idx]))

	check(err)

	return c
}

func getLogData() (int, int){
	absPath := getAbsolutePath("log.txt")

	if _, err := os.Stat(absPath); err == nil {
		var valuesInt []int
		rawText, err := ioutil.ReadFile(absPath)

		check(err)

		text := string(rawText)
		valuesStr := regexp.MustCompile(`\d+`).FindAllString(text, -1)

		for i := 0; i < len(valuesStr); i++ {
			num, _ := strconv.Atoi(valuesStr[i])
			valuesInt = append(valuesInt, num)
		}

		return valuesInt[0], valuesInt[1]
	}

	return 0, 0
}

func writeLogFile(chunkNum int, digit int){
	str := strconv.Itoa(chunkNum) + "," + strconv.Itoa(digit)
	bt := []byte(str)
	err := ioutil.WriteFile("log.txt", bt, 0644)
	check(err)
}

func check(e error){
	if (e != nil) {
		panic(e)
	}
}