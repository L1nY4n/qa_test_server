#!/usr/bin/sh

app=app=qa_test_server
p=`ps -ef |grep $app |grep -v grep |awk '{print $2}'`
if [ -n "$p" ] ; then
    echo "kill ${app}[$p]"
    kill -9  $p
else
    echo "${app} is not running."
fi