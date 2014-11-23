#!/bin/bash
./bazooka_before_install.sh
rc=$?
if [[ $rc != 0 ]] ; then
    exit $rc
fi
./bazooka_install.sh
rc=$?
if [[ $rc != 0 ]] ; then
    exit $rc
fi
./bazooka_before_script.sh
rc=$?
if [[ $rc != 0 ]] ; then
    exit $rc
fi
./bazooka_script.sh
if [[ $? != 0 ]] ; then
  ./bazooka_after_failure.sh
else
  ./bazooka_after_success.sh
fi
./bazooka_after_script.sh
