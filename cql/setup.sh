#!/bin/bash
host=localhost
username=
password=

files='_cleanup.cql generator.cql domain.cql session.cql event-store.cql trace.cql task.cql data.cql'


for file in $files
do
    echo ${file}
    #cqlsh ${host} -u ${username} -p ${password} -f ${file}
    cqlsh -f ${file}
done