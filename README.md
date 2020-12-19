# dirsplitter
Split large directories into parts of a specified maximum size

Check out my nim port for 80% smaller binary: https://github.com/jinyus/nim_dirsplitter

How to build:  
-Clone this git repo  
-cd into directory and run "go build" or "go install" to add executable to your path

Or download the prebuild binary from: https://github.com/jinyus/dirsplitter/releases


Usage of dirsplitter:  
  -dir string  
        &nbsp;&nbsp;&nbsp;&nbsp;Target Directory (default ".")  
  -max float  
        &nbsp;&nbsp;&nbsp;&nbsp;Max part size in GB (default 5)  
        
```
example:

dirsplitter -dir ./mylarge2GBdirectory -max "0.5"
NB: decimals has to be wrapped in quotes("")

This will yield the following directory structure:

📂mylarge2GBdirectory
 |- 📂part1
 |- 📂part2
 |- 📂part3
 |- 📂part4

with each part being a maximum of 500MB in size.
```
