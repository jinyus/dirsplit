# dirsplitter
Split large directories into parts of a specified maximum size

How to build:  
-Clone this git repo  
-cd into directory and run "go build" or "go install"


Usage of dirsplitter:  
  -folder string  
        &nbsp;&nbsp;&nbsp;&nbsp;Target folder (default ".")  
  -max float  
        &nbsp;&nbsp;&nbsp;&nbsp;Max folder size in GB (default 5)  
        
```
example
dirsplitter -folder ./mylargefolder -m 2
```
