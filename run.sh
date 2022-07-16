#!/usr/bin/sh

app=qa_test_server
processCount=`ps -ef | grep ${app} | grep -v grep|grep -v cd | wc -l`
if [ ${processCount} -eq 0 ] ; then
        echo "${app} is starting"
	      nohup ./${app} 2>&1 > /dev/null &
else
        echo "${app} is running."
fi