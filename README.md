# gofreedom


You know what I mean ...

Status
======

It is *not* usable yet and still under early development.

Build and Install
=====
    First you should config your golang environment
	
	cd  $GOPATH/src
    git clone https://github.com/oikomi/gofreedom.git
    cd gofreedom
    make
	
Usage
=====
	localserver ---- remoteserver(a vps can anti GFW)
	
	./bin/gofreedom ./config/httpproxy.json   (run in localserver)
	./bin/gofreedom ./config/tcpproxy.json   (run in remoteserver)
	
	config your web browser with httpproxy
	
Todo
======
- encrypt connetion(support aes des...)
- DNS pollution
- ...

Copyright & License
===================

Copyright 2014 Hong Miao. All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.