#!/usr/bin/env bash

virtualenv .venv
source .venv/bin/activate

pip install -r plugin/pyFile/requirements.txt