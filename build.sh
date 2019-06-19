#!/bin/sh

rootPwd=$(pwd)

# XXX This still reports errors because of missing github.com/modell-aachen/gologger
go build

success=0
for cmd in $(ls $rootPwd/cmd -1); do
    cd $rootPwd/cmd/$cmd && go build
    buildOk=$?
    if [ $buildOk != 0 ]; then
        success=$buildOk
    fi
done

exit $success
