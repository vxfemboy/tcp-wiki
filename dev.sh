#/bin/bash
# this sets up super annoying shit like hard symlinks and whatever else
# Eventually front end will be in a seperate branch for production
# but for now we can manage by editing via hardlinks. 
# This sets up your local readme as the main page and the files in assets as public
cd src && go run .
rm ../data/assets/*
ln ../assets/_layout.html ../data/assets/_layout.html
ln ../assets/main.css ../data/assets/main.css
rm ../data/README.md
ln ../README.md ../data/README.md
echo "Developer setup ready!"
echo "Goto: http://127.0.0.1:8080"

