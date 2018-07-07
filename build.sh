#!/bin/bash

dirPath=$(cd "$(dirname "$0")";pwd)

cd $dirPath/web/static
echo " cd web/static \n"

if [ $? != 0 ];then
 	echo "cd folder web/static failed \n"
 	exit
fi


{ # your 'try' block
   npm run build
} || { # your 'catch' block
    echo " npm run build failed \n"
}

echo " cd web/router \n"
cd $dirPath/web/router

if [ $? != 0 ];then
 	echo "cd folder  web/router failed \n"
 	exit
fi

echo " exec packr \n"

packr

if [ $? != 0 ];then
 	echo "packr failed \n"
 	exit
fi

cd $dirPath/_example
echo " cd _example \n"
if [ $? != 0 ];then
 	echo "cd folder  _example failed \n"
 	exit
fi

echo " exec go build main \n"
go build main.go
