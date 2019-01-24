#!/bin/bash

[ ! -f flower ] && ./build.sh

./flower -stderrthreshold=0 -v=2 -f=config.json -logtostderr -web=true -cron=false -worker=false
