#!/bin/bash
mongoimport --db PDDB --collection PowerDemand --type csv --headerline --file /files/power-demand.csv