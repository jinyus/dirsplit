# dirsplitter
Split large directories into parts of a specified maximum size

How to build:
-Clone this git repo
-cd into directory and run "go build" or "go install"


Usage of dirsplitter:
  -folder string
        Target folder (default ".")
  -max float
        Max folder size in GB (default 5)
        
```
example
dirsplitter -folder ./mylargefolder -m 2
```
