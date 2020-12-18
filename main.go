package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
)

func main() {
	//use float to allow specifying sizes < 1GB
	var maxSize float64
	var folder string
	flag.Float64Var(&maxSize, "max", 5, "Max folder size in GB")
	flag.StringVar(&folder, "folder", ".", "Target folder")

	flag.Parse()

	confirmOperation(fmt.Sprintf(`Splitting "%s" into %.3fGB parts.`, folder, maxSize))
	fmt.Printf("Slitting Directory\n\n")

	const GBMultiple = 1024 * 1024 * 1024
	tracker := map[int]float64{}
	currentPart := 1
	filesMoved := 0

	maxFileSize := maxSize * GBMultiple

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("could not access file : %q: %v\n", path, err)
			return err
		}

		if info.IsDir() && path != folder {
			// only proccess files in target folder. ie: depth of 1
			return filepath.SkipDir
		} else if path == folder {
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

		err = moveFile(path, info.Name(), folder, currentPart)
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
		log.Fatal(err)
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

func moveFile(fullPath, filename, dstFolder string, part int) error {
	partFolder := fmt.Sprintf("part%d", part)
	_ = os.Mkdir(path.Join(dstFolder, partFolder), os.ModePerm)
	finalDst := path.Join(dstFolder, partFolder, filename)
	return os.Rename(fullPath, finalDst)
}
