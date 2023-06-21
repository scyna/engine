#!/bin/bash
host=localhost
username=
password=

files='_cleanup.cql generator.cql domain.cql gateway.cql session.cql trace.cql task.cql data.cql gateway.cql'


for file in $files
do
    echo ${file}
    #cqlsh ${host} -u ${username} -p ${password} -f ${file}
    cqlsh -f ${file}
done