#!/bin/bash

echo 'Build go project'
go build -o kubectl-refresh

echo 'Remove old'
rm -rf /usr/local/bin/kubectl-refresh

echo 'Install new'
cp kubectl-refresh  /usr/local/bin