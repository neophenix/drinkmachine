# DrinkMachine

DrinkMachine is a raspberry pi powered bartender.  Using the provided [parts list](hardware/README.md) and wiring diagrams you can load this onto your pi and have it serving you drinks.

## Installation

Start with a lite [Raspbian](https://www.raspberrypi.org/downloads/raspbian/) install.  Then install a few packages to make life easier
```
sudo apt install git sqlite3
````

Next install golang, currently I use 1.10.3, to get that, and setup your paths, do:
```
wget https://storage.googleapis.com/golang/go1.10.3.linux-armv6l.tar.gz
sudo tar -C /usr/local -xzf go1.10.3.linux-armv6l.tar.gz
echo "export PATH=\$PATH:/usr/local/go/bin:~/go/bin" >> ~/.profile
echo "export GOPATH=\$HOME/go" >> ~/.profile
```

Now either `go get github.com/neophenix/drinkmachine` or otherwise copy this to your pi.

Setup the sqlite DB you'll need by using the .schema file in the db dir
```
sqlite3 drinkmachine.db < db/drinkmachine.schema
```
If you modified the wiring at all, you'll want to go into the DB and modify the pumps table, setting the correct GPIO pins for each pump.  See [db/README.md](db/README.md) for details on the db tables.

Last, build the binary and start the process
```
go build cmd/drinkmachine/drinkmachine.go
sudo ./drinkmachine
```
Note: sudo is required because the GPIO lib needs it and by default we want to use port 80

## Options

DrinkMachine has only a few command line flags to keep things simple
```
    -port       (default:80) which port to bind to
    -db         (default:drinkmachine.db) location of the sqlite3 DB
    -cache_templates (default:true) more for dev work, setting to false make the service read the web templates off disk each request
```

## Usage

Everything is done via the web interface (port 80 by default).  The main page lets you select a drink and pour it, then there are 3 "admin" pages

 * Manage Ingredients is the first place you should go, add all the alcohols, etc that you will want to use
 * Manage Pumps lets you select which ingredient is hooked up to each pump, it also lets you set the flow rate of the pump in ml/min.  The ones I used are rated for 100 ml/min and that is what the default is.
 * Manage Drinks is where you go to add all the cocktails you will want to make.  You give them a name, select an amount, units (the 3 common ones I know are defined, more can be added, PR welcome), and ingredient, and keep adding more ingredients as needed.  Then add some notes on finishing the drink or anything else you couldn't define elsewhere and save.
