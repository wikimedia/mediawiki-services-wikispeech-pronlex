## Docker installation

WORK IN PROGRESS

1. Build the Docker image

    1. from Dockerfile URL:

        `$ docker build https://raw.githubusercontent.com/stts-se/pronlex/master/Dockerfile -t stts-lexserver-local`   

<!---   $ docker build --build-arg USER=$USER https://raw.githubusercontent.com/stts-se/pronlex/master/Dockerfile -t stts-lexserver-local	 --->

    2. from local Dockerfile:

        `$ docker build $GOPATH/src/github.com/stts-se/pronlex -t stts-lexserver-local`

<!---       $ docker build --build-arg USER=$USER $GOPATH/src/github.com/stts-se/pronlex -t stts-lexserver-local --->


    Insert the `--no-cache` switch after the `build` tag if you encounter caching issues (updated git repos, etc).


2. Run the docker app


   1. Setup the server 

      `$ docker_run.sh -a <APPDIR> setup`


   2. Import lexicon files (optional)

      `$ docker_run.sh -a <APPDIR> import_all`


   3. Run lex server

      `$ docker_run.sh -a <APPDIR>`


   You can also investigate the server environment using `bash`:   

   `$ docker_run.sh -a <APPDIR> bash`
  

