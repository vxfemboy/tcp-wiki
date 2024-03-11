
# TCP-WIKI
[![Go Report Card](https://goreportcard.com/badge/github.com/vxfemboy/tcp-wiki)](https://goreportcard.com/report/github.com/vxfemboy/tcp-wiki)

![screenshot](examples/tcp-wiki.png "TCP-WIKI")

### What is TCP-WIKI ? 

TCP.WIKI is a secure, minimal, and verifiable wiki platform designed for projects, code, courses, documents, articles, blogs, tutorials, and more.

### Project Goals

The aim is to provide a secure, minimal, and easily verifiable wiki environment that supports a wide range of content types, from technical documentation, to educational materials, to blogs, and more.

#### Features:

* Full support for pages in Markdown and HTML

* Pull from GIT for live updates or usage of local directory

* No Javascript required

## Setup

First clone this repository:
```bash
git clone https://github.com/vxfemboy/tcp-wiki.git
```
Then you have to cd into the repo's folder and run/compile:
```bash
cd tcp-wiki
go run ./src
```
Then you goto your browser and visit: http://127.0.0.1:8080/

## Want to use with your own data?

All you have to do is modify the following lines in the `config.toml` file:

```toml
[Git]
UseGit = true # Set to false to use LocalPath
RepoURL = "https://github.com/vxfemboy/tcp-wiki.git" # Your Repo Here
Branch = "main" # Your Repo Branch Here
LocalPath = "data" # Directory to clone the git repo too
```

Change the `RepoURL` line `https://github.com/vxfemboy/tcp-wiki.git` to your repo link,
change `main` to your specific repo's branch and you should be good to go!

#### Want to use a local directory other then git repo?

To do this you just need to set `UseGit` to `false` and set your directory in config.toml

```toml
[Git]
UseGit = false # Set this to false 
RepoURL = "" # Ignored
Branch = "" # Ignored
LocalPath = "/home/crazy/blog" # The directory of your project
```
make sure to also set `LocalPath` to the directory of your project

#### Want to use your own theme/layout?

Have a look at the `assets/` directory for the template `_layout.html` this is the main file that will be used

using [tailwindscss](https://tailwindcss.com/) and [daisyui](https://daisyui.com/) as a css library 

to setup this for live modifications do:
```
cd sass
npm install
npx tailwindcss -i ./tail.css -o ../assets/main.css --watch
```
and to minify your css for production usage:
```
npx tailwindcss -o ../assets/main.css --minify
```
> ### ❤️ Feel free to commit, leave suggestions/ideas, issues, or really anything ❤️

## TODO

- [x] config file
- [ ] Webhook support for auto pull on push/update of the git repo
- [x] Git Branch support
- [ ] add a star/upvote/like feature for pages
- [x] edit/version tracker 
    - [x] Author 
    - [x] last edited
    - [x] last editor/commit - maybe working
    - [ ] PGP Signed & Verification
- [ ] pgp signed intergration
- [x] comments using bitcask - generated in comments.db/
    - [ ] verification - no login pgp
    - [ ] captcha
    - [ ] sub rating system
    - [ ] sort by date etc
    - [ ] reply to replies
    - [ ] set security controls per page
    - [ ] auto refresh on post
- [x] dynamically generated links for all avaiable pages
    - [x] sitemap - kinda (needs limited)
    - [x] working pages
    - [x] image support
