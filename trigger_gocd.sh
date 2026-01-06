#!/bin/sh

curl -X POST \
  -H "Accept: application/vnd.go.cd.v1+json" \
  -H "X-GoCD-Confirm: true" \
  http://localhost:8153/go/api/pipelines/go-app/schedule

