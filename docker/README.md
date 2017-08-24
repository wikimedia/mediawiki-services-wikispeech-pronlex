## Docker installation

The `<DOCKERTAG>` variable below should be set to `stts-lexserver-local`.

1. Build the Docker image

    1. from Dockerfile URL:

        `$ docker build https://raw.githubusercontent.com/stts-se/pronlex/master/Dockerfile -t <DOCKERTAG>`   

    2. from local Dockerfile:

        `$ docker build $GOPATH/src/github.com/stts-se/pronlex -t stts-lexserver-local`

    Insert the `--no-cache` switch after the `build` tag if you encounter caching issues (e.g., with updated git repos, etc).


2. Run the docker app


   1. Setup the server 

      `$ docker_run.sh -a <APPDIR> -t <DOCKERTAG> setup`

      Set up the server's required file structure in the specified `<APPDIR>`
      

   2. Import lexicon files (optional)

      `$ docker_run.sh -a <APPDIR> -t <DOCKERTAG> import_all`

      Imports lexicon data for Swedish, Norwegian, US English and a small test file for Arabic.


   3. Run lex server

      `$ docker_run.sh -a <APPDIR> -t <DOCKERTAG>`


   You can also investigate the server environment using `bash`:   

   `$ docker_run.sh -a <APPDIR> -t <DOCKERTAG> bash`
  

Server data files and databases are saved in the folder `<APPDIR>`. Please note that this folder will be owned by `root`. If this is a problem, make sure you change the ownership and/or permissions to whatever is best for your environmemnt.


<!-- to pass on system user to the docker environment:
<!---   $ docker build --build-arg USER=$USER https://raw.githubusercontent.com/stts-se/pronlex/master/Dockerfile -t stts-lexserver-local	 --->

<!---       $ docker build --build-arg USER=$USER $GOPATH/src/github.com/stts-se/pronlex -t stts-lexserver-local --->


