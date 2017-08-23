## Docker installation

WORK IN PROGRESS

1. Build

    1. from Dockerfile URL:

        `$ docker build https://raw.githubusercontent.com/stts-se/pronlex/master/docker/Dockerfile -t sttsse/lexserver`   

    2. from local Dockerfile:

        `$ docker build $GOPATH/src/github.com/stts-se/pronlex/docker -t sttsse/lexserver`

    Insert the `--no-cache` switch after the `build` tag if you encounter caching issues (updated git repos, etc).


2. Create the application folder:

   `$ mkdir <APPDIR>`


3. Setup server

   ``$ docker run -u `stat -c "%u:%g" <APPDIR>` -v </FULL/PATH/TO/APPDIR>:/go/appdir -p 8787:8787 -it sttsse/lexserver setup``


4. Import lexicon files (optional)

   ``$ docker run -u `stat -c "%u:%g" <APPDIR>` -v </FULL/PATH/TO/APPDIR>:/go/appdir -p 8787:8787 -it sttsse/lexserver import_lex``


5. Run lex server

   ``$ docker run -u `stat -c "%u:%g" <APPDIR>` -v </FULL/PATH/TO/APPDIR>:/go/appdir -p 8787:8787 -it sttsse/lexserver``



You can also investigate the server environment using `bash`:

``$ docker run -u `stat -c "%u:%g" <APPDIR>` -v </FULL/PATH/TO/APPDIRAPPDIR>:/go/appdir -p 8787:8787 -it sttsse/lexserver bash``

