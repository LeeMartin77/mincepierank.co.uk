#!/bin/bash

# This is for situations where the automounting fails - probably from being on kind/linux


kubectl cp -n mincepierank-development testdata/images mincepierank-imgprssr:/
