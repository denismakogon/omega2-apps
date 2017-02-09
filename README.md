Omega Onion 2 + Go 1.8 = <3
===========================

Recent (or upcoming release) of Go 1.8 brings a lot features like SO plugins and major improvements.
But i'm so hyped by hearing that GO 1.8 now supports MIPS32LE cross compiling!

Before Go 1.8
-------------

I'd like to say just a few words: `gccgo` from sources, custom OpenWRT builds with gccgo toolchains
IMHO, too much problems with such small devices like IoT thing.

With Go 1.8
-----------

With [new cross-compiling features](https://beta.golang.org/doc/go1.8#ports) of Go 1.8rc3 i was able build a [Gorilla Mux](https://github.com/gorilla/mux) web server and run it successfully on Onion Omega 2

Here's what you need
```bash
$ export GOPATH=~/go
$ go get -u github.com/gorilla/mux
```    

In order to build our [simple http server](mux_server.go) you need to use following command:
```bash
$ GOOS=linux GOARCH=mipsle go build -compiler gc mux_server.go
```
_**Note** the *cross-compilation* options: ` GOOS=linux GOARCH=mipsle ` that aim for *linux* and *MIPS Little Endian* (Omega 2's architecture)._

Boom! Now you have executable binary app, copy it over WiFi to you Omega2

```bash
$ scp mux_server root@192.168.3.1:/root
```

Take a look at [Omega 2 Getting started guide](https://wiki.onion.io/get-started) for more information how to connect to your device.

Here's an example for [Hello World app](hello_world.go):

```bash
BusyBox v1.25.1 () built-in shell (ash)

   ____       _             ____
  / __ \___  (_)__  ___    / __ \__ _  ___ ___ ____ _
 / /_/ / _ \/ / _ \/ _ \  / /_/ /  ' \/ -_) _ `/ _ `/
 \____/_//_/_/\___/_//_/  \____/_/_/_/\__/\_, /\_,_/
 W H A T  W I L L  Y O U  I N V E N T ? /___/
 -----------------------------------------------------
   Ω-ware: 0.1.9 b149
 -----------------------------------------------------
root@IronOmega:~# ls -la
drwxr-xr-x    1 root     root           0 Feb  5 21:07 .
drwxr-xr-x    1 root     root           0 Feb  5 20:02 ..
-rwxr-xr-x    1 root     root        3.7M Feb  5 21:08 mux_server
-rwxr-xr-x    1 root     root        1.0M Feb  5 22:20 hello_world

root@IronOmega:~# ./hello_world 
Hello linux/mipsle
```


Negative effect
---------------

[mux_server](mux_server.go) file with 25 LOC has size of 539 bytes, but its binary file weights something around 5.5Mb.
Onion Omega 2 chip comes with 14Mb storage available, so there's no so much space for a good party, unfortunately.

Positive effect
---------------

Thanks to quite impressive Go capabilities to shrink its binaries and make them more affordable for IoT devices with limited storage. 

```bash
$ GOOS=linux GOARCH=mipsle go build -ldflags "-s -w" -compiler gc mux_server.go
```

Boom! And your binary gets lighter up to 2Mb (almost half size of initial).


Conclusions
-----------

Go 1.8 made great progress in its cross-compiling features along with capability to shrink compiled binaries into affordable size.
Go made a good step forward into IoT world as de-facto standard development language for MIPS32 devices.
