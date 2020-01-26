
#!/bin/bash
#
# remove all images matching a path by their digest
#
repository=localhost:32000
path=""
mode="digest"

cleanPath()
{
  echo "Checking repository: $repository"
  echo "for images matching: $path'*'"
  fileMatch=$(docker image ls --all --digests $repository/$path"*"); echo "$fileMatch"
  echo "Found the following digests"
  read -p "Confirm deletion by entering 'yes': " yesNo
  yesNo=`echo ${yesNo^^}`

  if [ $yesNo == "YES" ]
  then
    toDelete=`docker image ls --all --digests $repository/$path"*" | awk '{ if(NR>1) print $4 }'`
    for imageID in $toDelete
    do
      if [[ $imageID != *"<none>"* ]]; then
        echo "deleting: docker image rm -f $imageID"
        docker image rm -f $imageID
      fi
    done
  fi
}

usage()
{
  echo "usage: filmsRepositoryClean.sh -p /dir/dir/images"
  echo "    Retrieves a list of image digests"
  echo "    append '*' on the end of the path given"
}

## Main
while [ "$1" != "" ]; do
  case $1 in
    -p | --path ) shift
      path=$1
      ;;
  -h | --help ) usage
    exit
    ;;
  * ) shift
    ;;
  esac
done

# call all
cleanPath
