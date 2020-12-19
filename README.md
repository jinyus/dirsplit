# dirsplitter
Split large directories into parts of a specified maximum size

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

dirsplitter -dir ./mylarge8GBdirectory -max 0.5

This will yield the following directory structure:

📂mylarge2GBdirectory
 |- 📂part1
 |- 📂part2
 |- 📂part3
 |- 📂part4

with each part being a maximum of 500MB in size.
```
