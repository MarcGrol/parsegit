#!/bin/sh

sudo easy_install pip
sudo pip install --upgrade git2json

cd $GOPATH/src/github.com/Duxxie/platform
git2json | python -m json.tool > git_history.json

cd $GOPATH/src/github.com/MarcGrol/parsegit

#go get github.com/ChimeraCoder/gojson/...
#gojson -input git_history.json > model.go


go install &&  parsegit -filename=git_history.json