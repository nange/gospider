#!/bin/bash

dirPath=$(cd "$(dirname "$0")";pwd)

cd $dirPath/web/static

if [ $? != 0 ];then
 	echo "cd folder web/static failed \n"
 	exit
fi

#catch this command error 
{
	if [ ! -d "./node_modules" ]; then
		echo "not found node_modules; will exec npm install \n"
		npm install
	fi
	npm run build:prod
} || {
	echo "exec npm run build failed\n"
}

cd $dirPath/web/router

if [ $? != 0 ];then
 	echo "cd folder web/router failed\n"
 	exit
fi

packr

if [ $? != 0 ];then
 	echo "packr failed \n"
 	exit
fi

cd $dirPath/_example

if [ $? != 0 ];then
 	echo "cd folder _example failed\n"
 	exit
fi

go build

if [ $? != 0 ];then
 	echo "exec go build main.go failed\n"
 	exit
fi
echo "success!"

