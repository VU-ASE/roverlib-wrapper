import os
import time
import sys

#
# This is a basic script to confirm that the ServiceWrapper works for you. It shows how to flush stdout so that output can be shown in the console.
#
# Run this script as follows:
# - make start ruanrgs="'python3 test.py'"
# OR (after build)
# ./bin/mod-ServiceWrapper "python3 test.py"
#

# Get env variable
# service_name = os.environ.get('ASE_SW_ServiceName')
while True:
    for key, value in os.environ.items():
        if "ASE" in key:
            print(f"{key}: {value}")
    sys.stdout.flush()
    time.sleep(2)