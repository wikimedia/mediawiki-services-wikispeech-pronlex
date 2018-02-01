folder=vendor/js
outputfile=dependencies.txt
cmd=`basename $0`

git grep "\<script.*javascript.*\.js" | egrep -v "\<.-- "| sed 's/<script.*src="//' | sed 's/">.*//' >| $folder/$outputfile && echo "[$cmd] Javascript dependencies written to $folder/$outputfile"

cd $folder
cat $outputfile | egrep http | sed 's/.*\b\(http\)/\1/' | sort -u | xargs wget

ndeps=`cat $outputfile | egrep http | sed 's/.*\b\(http\)/\1/' | sort -u | wc -l`

echo "[$cmd] $ndeps downloaded to $folder"
