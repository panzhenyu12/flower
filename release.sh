#!/bin/bash

./build.sh
echo -e \\n

### flower
version=`cat VERSION`

mkdir -p flower-$version
#rm flower-$version.tar.gz || true
cp flower flower-$version/
cp config.json flower-$version/
cp config.yml flower-$version/
cp ChangeLog.md flower-$version/
cp build.sql flower-$version/deepface.sql
cp run.sh flower-$version/
#cp -rf file/ flower-$version/file/
#cp _flower.postman_collection.json flower-$version/
cp VERSION flower-$version/

rm -rf ~/Desktop/release/flower || true
rm ~/Desktop/release/flower-$version.tar.gz || true
mkdir ~/Desktop/release/flower
cp -rf flower-$version ~/Desktop/release/flower
rm -rf flower-$version
cd ~/Desktop/release/flower
ln -s flower-$version latest
chmod 777 latest
cd ../

sleep 1s
tar -zcvf flower-$version.tar.gz flower
#rm -rf flower-$version
