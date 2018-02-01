outputfile=vendor/js_dependencies.txt
git grep "\<script.*javascript.*\.js" | egrep -v "\<.-- "| sed 's/<script.*src="//' | sed 's/">.*//' >| $outputfile && echo "Javascript dependencies written to $outputfile"
