outdir=`dirname $0 | xargs realpath`
gitroot=`echo $outdir/../.. | xargs realpath`

outfile=dependencies.txt
cmd=`basename $0`

echo $gitroot

cd $gitroot && git grep "\<script.*javascript.*\.js" | egrep -v "\<.-- "| sed 's/<script.*src="//' | sed 's/">.*//' >| $outdir/$outfile && echo "[$cmd] Javascript dependencies written to $outdir/$outfile"

cd $outdir
cat $outfile | egrep http | sed 's/.*\b\(http\)/\1/' | sort -u | xargs wget -c

ndeps=`cat $outfile | egrep http | sed 's/.*\b\(http\)/\1/' | sort -u | wc -l`

echo "[$cmd] $ndeps downloaded to $folder"

cd -
