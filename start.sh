#!/bin/bash

nohup ./zion-tool tps --config=config/local.json --num=30 --period=480h30m --txn=10 --inc=10000 >> stable.log 2>&1 &
