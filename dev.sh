#/bin/bash
# this sets up super annoying shit like hard symlinks and whatever else
# but for now we can manage by editing via hardlinks. 
# This sets up your local readme as the main page and the files in assets as public

# Clone the repository
echo "Press Control+C when prompted"
go run ./src



# Set up hard links 
# !!! for main branch only !!!
rm data/assets/*
ln assets/_layout.html data/assets/_layout.html
ln assets/main.css data/assets/main.css
rm data/README.md
ln README.md data/README.md

echo "Developer setup ready!"
echo "Go ahead and run go run src/"
echo "And Go to: http://127.0.0.1:8080"
