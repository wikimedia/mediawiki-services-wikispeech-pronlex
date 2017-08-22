## WORK IN PROGRESS

# 1. Clone the source code

#     $ mkdir -p $GOPATH/src/github.com/stts-se
#     $ cd $GOPATH/src/github.com/stts-se
#     stts-se$ git clone https://github.com/stts-se/pronlex.git


# 2. Download dependencies
    
#     $ cd $GOPATH/src/github.com/stts-se/pronlex
#     pronlex$ go get ./...

# 3. Clone the lexdata repository
    
#      $ mkdir -p ~/gitrepos  
#      $ cd ~/gitrepos  
#      gitrepos$ git clone https://github.com/stts-se/lexdata.git


# 4. Prepare symbol sets and symbol set mappers/converters
    
#      $ cd $GOPATH/src/github.com/stts-se/pronlex/lexserver
#      lexserver$ mkdir symbol_sets  
#      lexserver$ cp ~/gitrepos/lexdata/*/*/*.sym symbol_sets   
#      lexserver$ cp ~/gitrepos/lexdata/mappers.txt symbol_sets  
#      lexserver$ cp ~/gitrepos/lexdata/converters/*.cnv symbol_sets  


