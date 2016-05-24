#!/bin/bash
rm $GOPATH/bin/storage*
rm $GOPATH/bin/client
go install github.com/ManikDV/storage/client
go install github.com/ManikDV/storage/storage