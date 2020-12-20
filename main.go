package main

import (
	"flag"
	"fmt"
	cop "github.com/jinyus/confirmop"
	"os"
	"path"
	"path/filepath"
)

func main() {
	//use float to allow specifying sizes < 1GB
	var maxSize float64
	var dir string
	var showTarCommand bool
	var outPrefix string
	flag.Float64Var(&maxSize, "max", 5, "Max part size in GB")
	flag.StringVar(&dir, "dir", ".", "Target directory")
	flag.BoolVar(&showTarCommand, "show-cmd", false, "Show tar command to compress each directory")
	flag.StringVar(&outPrefix, "out-prefix", "", "Prefix for output files of the tar command. -show-cmd must be specified. eg: myprefix.part1.tar")

	flag.Parse()

	userChoice := cop.ConfirmOperation(fmt.Sprintf(`Splitting "%s" into %.3fGB parts.`, dir, maxSize), "continue", false)

	if !userChoice {
		os.Exit(0)
	}
	fmt.Printf("Splitting Directory\n\n")
	if outPrefix != "" {
		outPrefix += "."
	}

	const GBMultiple = 1024 * 1024 * 1024
	tracker := map[int]float64{}
	currentPart := 1
	filesMoved := 0
	failedOps := 0

	maxFileSize := maxSize * GBMultiple

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("could not access file : %q: %v\n", path, err)
			return err
		}

		if info.IsDir() && path != dir {
			// only proccess files in target dir. ie: depth of 1
			return filepath.SkipDir
		} else if path == dir {
			return nil
		}

		fileSize, err := getFileSize(path)
		if err != nil {
			fmt.Printf("could not get file size : %v\n", err)
		}

		decrementIfFailed := false
		if tracker[currentPart]+fileSize > maxFileSize && tracker[currentPart] > 0 {
			currentPart++
			decrementIfFailed = true
		}
		tracker[currentPart] += fileSize

		err = moveFile(path, info.Name(), dir, currentPart)
		if err != nil {
			fmt.Printf("could not move file : %v\n", err)
			tracker[currentPart] -= fileSize
			failedOps++
			if decrementIfFailed {
				currentPart--
			}

		}
		filesMoved++
		return nil
	})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if filesMoved == 0 {
		currentPart = 0
	}
	fmt.Printf("Done:\nParts created: %d\nFiles moved: %d\nFailed Operations: %d\n", currentPart, filesMoved, failedOps)

	if currentPart > 0 && showTarCommand {
		if currentPart == 1 {
			fmt.Printf(`Tar Command : tar -cf "%spart1.tar" "part1"; done`, outPrefix)
		} else {
			fmt.Printf(`Tar Command : for n in {1..%d}; do tar -cf "%spart$n.tar" "part$n"; done`, currentPart, outPrefix)
		}
	}

}

func getFileSize(filename string) (float64, error) {
	fi, err := os.Stat(filename)
	if err != nil {
		return 0, err
	}
	return float64(fi.Size()), nil
}

func moveFile(fullPath, filename, dstDir string, part int) error {
	partDir := fmt.Sprintf("part%d", part)
	_ = os.Mkdir(path.Join(dstDir, partDir), os.ModePerm)
	finalDst := path.Join(dstDir, partDir, filename)
	return os.Rename(fullPath, finalDst)
}
