"use strict";

angular.module('bzk.project').controller('ProjectController', function($scope, BzkApi, $routeParams, $interval, $location) {
    var pId = $routeParams.pid;

    BzkApi.project.get(pId).success(function(project) {
        $scope.project = project;
    });

    $scope.refreshJobs = function() {
        BzkApi.project.jobs(pId).success(function(jobs) {
            $scope.jobs = jobs;
        });
    };

    $scope.refreshJobs();

    $scope.newJob = {
        reference: 'master'
    };

    $scope.startJob = function() {
        BzkApi.project.build($scope.project.id, $scope.newJob.reference).success(function() {
            $scope.refreshJobs();
            $scope.showNewJob = false;
        });
    };

    var refreshPromise = $interval($scope.refreshJobs, 5000);
    $scope.$on('$destroy', function() {
        $interval.cancel(refreshPromise);
    });
});