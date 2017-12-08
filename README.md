# Parse and process git history of any project

    go get github.com/MarcGrol/parsegit


## Extract history from project

Use git2json tool (see https://github.com/tarmstrong/git2json)

    $ cd <your git project>
    $ git2json | python -m json.tool > git_history.json

## Process the history

*Analyse committers*

    $ parsegit -filename=git_history.json -analyze-committers | sort -n -t';' -r -k2

    $ parsegit -filename=git_history.json -analyze-files | sort -n -t';' -r -k2
     
     

    
