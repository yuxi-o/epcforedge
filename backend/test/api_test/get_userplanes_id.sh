
#! /bin/sh
#########################################################
# <COPYRIGHT_TAG>
#########################################################
setup_dir=${PWD}

set -e

curl -v --cacert mec.crt https://mec.local:8080/userplanes/5
#curl -v --cacert mec.crt https://mec.local:8080/userplanes/123