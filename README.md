H0wdy!!!

feel free to commit, leave suggestions/ideas, issues, or really anything <3

# What is TCP.WIKI ? 

## Project Goals
secure and verifiable wiki for projects, code, courses, documents, articles, tutorials, and more

## Setup
**For a normal user you can follow this process:**

First clone the repo:
```bash
git clone https://git.tcp.direct/S4D/tcp-wiki.git
```
Then you have to cd into the repo's folder and run/compile:
```bash
cd tcp-wiki/src
go run .
```
Then you goto your browser and visit: http://127.0.0.1:8080/

**For a develeper setup you can follow this process:**

First clone the repo:
```bash
git clone ssh://git@git.tcp.direct:2222/S4D/tcp-wiki.git
```
Then cd and run dev.sh
```bash
cd tcp-wiki
bash dev.sh
```
Then you just have to execute this to run the server:
```go
cd src && go run .
```
Then you goto your browser and visit: http://127.0.0.1:8080/

This method just adds in some handy symlinks for development purposes

## Want to use with your own repo?

All you have to do is modify the following lines in the src/main.go file:
```go
const repoURL = "https://git.tcp.direct/S4D/tcp-wiki.git"
```
Change `https://git.tcp.direct/S4D/tcp-wiki.git` to your repo link, and:
```go
const repoBRANCH = "main"
```
change `main` to your specific repo's branch and you should be good to go!

## TODO

- [ ] MANY FUCKING THINGS
- [ ] Webhook support for auto pull on push/update of the git repo
- [x] Git Branch support
- [ ] add a star/upvote/like feature for pages
- [ ] edit tracker 
    - [ ] Author 
    - [ ] last edited
    - [ ] last editor/commit
    - [ ] PGP Signed & Verification
- [ ] pgp signed intergration
- [x] comments using bitcask - generated in comments.db/
    - [ ] verification - no login pgp
    - [ ] captcha
    - [ ] sub rating system
    - [ ] sort by date etc
    - [ ] reply to replies
    - [ ] set security controls per page
- [ ] dynamically generated links for all avaiable pages
    - [ ] sitemap 
    - [ ] anti robot shit here
    - [ ] acual working pages!?
- [ ] post quantum intergration and verification
- [ ] BUILD UP THAT MARKDOWN SUPPORT
- [ ] fix whatever i did to fuck up design/layout/css???