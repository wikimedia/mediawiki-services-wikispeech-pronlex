# How to merge in upstream

I.e. how to fetch changes from the remote upstream project repository and merge it to the gerrit repository.

One-off:
```
git remote add github https://github.com/stts-se/pronlex.git
```

At merge time:
```
git checkout master
git fetch --all
git pull
git push gerrit github/master:upstream
git checkout -b 202Y-MM-DD-merge-github
git merge gerrit/upstream
git review -R
```
