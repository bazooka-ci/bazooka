"use strict";

angular.module('bzk.project').controller('ProjectController', function($scope, BzkApi, $routeParams, $interval, $location){
	var pId = $routeParams.pid;

	BzkApi.project.get(pId).success(function(project){
		$scope.project = project;
	});

	$scope.refreshJobs = function() {
		BzkApi.project.jobs(pId).success(function(jobs){
			$scope.jobs = jobs;
		});
	};

	$scope.refreshJobs();

	$scope.newJob = {
		reference: 'master'
	};

	$scope.newJobVisible = function(s) {
		$scope.showNewJob = s;
	};

	$scope.startJob = function() {
		BzkApi.project.build($scope.project.id, $scope.newJob.reference).success(function(){
			$scope.refreshJobs();
			$scope.showNewJob = false;
		});
	};

	$scope.isSelected = function(j) {
		return j.id.indexOf($location.search().j)===0;
	};

	var refreshPromise = $interval($scope.refreshJobs, 5000);
	$scope.$on('$destroy', function() {
		$interval.cancel(refreshPromise);
	});
});
