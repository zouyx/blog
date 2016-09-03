echo "git checkout."
git co develop
echo "git pull."
git pull

APP_DIR=/appdata/deploy/blog
Project=blog
echo "make deploy folder."
sudo mkdir -p $APP_DIR

echo "build app."
go build -o "$Project"  "$Project"

echo "make outer."
rm -rf outer
mkdir outer

echo "copy file to outer."
cp -f run.sh outer/
cp -rf static outer/
cp -rf views outer/
cp -rf conf outer/
cp -rf $Project outer/
echo "deploy finished."

echo "delete file."
sudo rm -rf $APP_DIR/*
echo "copy file to deploy folder."
sudo cp -rf outer/* $APP_DIR