{{define "preflight.tpl"}}#!/bin/bash

defaultDirectory="."
prInitFile=".pr_git_initialize_file"
gitIgnore=".gitignore"
preSuccess=".pr_preflight_check"

makeInitFile()
{
	target=$prInitFile
	if [[ "$defaultDirectory" != "." ]]
	then
		target="$defaultDirectory/$prInitFile"
	fi
	cat << EOF > $target
Initial file for git commit generated by roadctl
feel free to delete
EOF
  echo "creating: $prInitFile"
}

initRepo()
{
	if [[ "$defaultDirectory" == "." ]]
	then
    (git log | grep fatal > /dev/null)
		if [ $? -eq 0 ]
		then
			exit 0
		fi
		(git init > /dev/null)
	else
    (cd $defaultDirectory;git log | grep fatal > /dev/null)
		if [ $? -eq 0 ]
		then
			exit 0
		fi
		(cd $defaultDirectory;git init > /dev/null)
	fi
}

addFile()
{
	if [[ "$defaultDirectory" == "." ]]
	then
		git add $prInitFile
	else
		(cd $defaultDirectory;git add $prInitFile)
	fi
  echo "adding: $prInitFile"
}

checkPreflight()
{
  if [[ -f "$preSuccess" ]]
  then
    exit 0
  fi
}

commit()
{
	if [[ "$defaultDirectory" == "." ]]
	then
		git commit -m "Initial commit of $prInitFile"
	else
		(cd $defaultDirectory;git commit -m "Initial commit of $prInitFile")
	fi
  echo "commit: $prInitFile"
}

checkUser()
{
	git config --global user.name > /dev/null

	if [ $? -ne 0 ]
	then
		read -p "Enter git user.name: " gituser
		git config --global user.name $gituser
	fi
}

checkEmail()
{
	git config --global user.email > /dev/null

	if [ $? -ne 0 ]
	then
		read -p "Enter git user.email: " gitemail
		git config --global user.email $gitemail
	fi
}

checkRepository()
{
	makeInitFile
	initRepo
	addFile
  commit
}

gitCheck()
{
	checkUser
	checkEmail
	echo '{"status":"success"}' > $preSuccess
	checkRepository
}

usage()
{
	echo "usage: preflight.sh -d"
	echo "		-d directory to test against, default is current "
}

## Main
while [ "$1" != "" ]; do
	case $1 in
		-d ) shift
			defaultDirectory=$1
			;;
	-h | --help ) usage
		exit
		;;
	* ) shift
		;;
	esac
done

checkPreflight
gitCheck
{{end}}
