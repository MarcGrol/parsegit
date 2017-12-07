# Parse and process git history of any project

    go get github.com/MarcGrol/parsegit


## Extract history from project

Use git2json tool (see https://github.com/tarmstrong/git2json)

    $ cd <your git project>
    $ git2json | python -m json.tool > git_history.json

## Process the history

    $ parsegit -filename=git_history.json -analyze-committers=true > committers.csv
     
     

    
