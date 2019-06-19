#!/usr/bin/env bash

set -euo pipefail

echo HELLO | sudo tee /opt/test.txt
