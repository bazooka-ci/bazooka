#!/bin/bash

if {{$.BzkBuildDir}}/bazooka_before_install.sh && \
   {{$.BzkBuildDir}}/bazooka_install.sh && \
   {{$.BzkBuildDir}}/bazooka_before_script.sh
then
	true
else
	exit 42
fi

{{$.BzkBuildDir}}/bazooka_script.sh
exitCode=$?


if [[ $exitCode == 0 ]]
then
  {{$.BzkBuildDir}}/bazooka_archive_success.sh
  {{$.BzkBuildDir}}/bazooka_after_success.sh
else
  {{$.BzkBuildDir}}/bazooka_archive_failure.sh
  {{$.BzkBuildDir}}/bazooka_after_failure.sh
fi

{{$.BzkBuildDir}}/bazooka_archive.sh
{{$.BzkBuildDir}}/bazooka_after_script.sh

exit $exitCode