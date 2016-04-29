#! /bin/bash
set -e

pdflatex(){
	/usr/local/texlive/2015/bin/universal-darwin/pdflatex "$@"
}

bibtex(){
	/usr/local/texlive/2015/bin/universal-darwin/bibtex "$@"
}

CHAPTER=$1

if [[  "$CHAPTER" == "" ]]; then
	pdflatex main.tex
	bibtex main
	pdflatex main.tex
	pdflatex main.tex
else
	pdflatex -jobname=chapBuild "\includeonly{$CHAPTER}\input{main.tex}"
	bibtex chapBuild
	pdflatex -jobname=chapBuild "\includeonly{$CHAPTER}\input{main.tex}"
	pdflatex -jobname=chapBuild "\includeonly{$CHAPTER}\input{main.tex}"
fi

