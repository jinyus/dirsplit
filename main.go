package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"path/filepath"
)

func main() {
	//use float to allow specifying sizes < 1GB
	var maxSize float64
	var dir string
	flag.Float64Var(&maxSize, "max", 5, "Max part size in GB")
	flag.StringVar(&dir, "dir", ".", "Target directory")

	flag.Parse()

	confirmOperation(fmt.Sprintf(`Splitting "%s" into %.3fGB parts.`, dir, maxSize))
	fmt.Printf("Slitting Directory\n\n")

	const GBMultiple = 1024 * 1024 * 1024
	tracker := map[int]float64{}
	currentPart := 1
	filesMoved := 0

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
	fmt.Printf("Success:\nParts created: %d\nFiles moved: %d\n", currentPart, filesMoved)

}

func confirmOperation(desc string) {
	var answer string
	fmt.Printf("%s \nconfirm? (y/n): ", desc)
	if _, err := fmt.Scanf("%s", &answer); err != nil {
		fmt.Printf("invalid answer : expected (y or n) got (%s) :\n%v", answer, err)
		os.Exit(1)
	} else if answer != "y" && answer != "n" {
		fmt.Printf("invalid answer : expected (y or n) got (%s)\n", answer)
		os.Exit(1)
	} else if answer != "y" {
		fmt.Println("Goodbye!")
		os.Exit(1)
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
