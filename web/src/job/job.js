"use strict";

angular.module('bzk.job').controller('JobLogsController', function($scope, JobResource, DateUtils, $routeParams, $timeout){
	var jId = $routeParams.jid;
	$scope.logger={};
	function loadLogs() {
		$scope.logger.job.prepare();

		JobResource.jobLog(jId).success(function(logs){
			$scope.logger.job.finish(logs);
		});
	}

	$timeout(loadLogs);

});

