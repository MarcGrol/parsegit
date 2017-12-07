# Parse and process git history of any project

    go get github.com/MarcGrol/parsegit


## Extract history from project

Use git2json tool

     git2json | python -m json.tool > git_history.json

## Process the history

    parsegit -filename=git_history.json -analyze-committers=true > committers.csv
     
     

    
