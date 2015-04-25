"use strict";

angular.module('bzk.home').controller('HomeController', function($scope, BzkApi, $interval){
	$scope.refreshJobs = function() {
		BzkApi.job.list().success(function(jobs){
			$scope.jobs = jobs;
		});
	};

	$scope.refreshJobs();

	var refreshPromise = $interval($scope.refreshJobs, 5000);
	$scope.$on('$destroy', function() {
		$interval.cancel(refreshPromise);
	});
});
