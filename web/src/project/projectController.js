"use strict";

angular.module('bzk.project').controller('ProjectController', function($scope, BzkApi, $routeParams, $interval, $location, growl) {
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
        parameters: []
    };

    $scope.addJobParam = function() {
        $scope.newJob.parameters.push({});
    };

    $scope.remJobParam = function(index) {
        $scope.newJob.parameters.splice(index, 1);
    };

    $scope.canStartJob = function() {
        return $scope.newJob.reference && _.every($scope.newJob.parameters, function(p) {
            return p.key && p.value;
        });
    };

    $scope.startJob = function() {
        var parameters = _.map($scope.newJob.parameters, function(p) {
            return p.key + "=" + p.value;
        });
        BzkApi.project.build($scope.project.id, $scope.newJob.reference, parameters).
        success(function(job) {
            growl.success('<p>Using SCM ref <strong>' + $scope.newJob.reference + '</strong><p>', {
                title: 'Job ' + $scope.project.name + '#' + job.number + ' started',
                ttl: 6000
            });
            $scope.newJob = {
                parameters: []
            };
            $scope.refreshJobs();
        });
    };

    var refreshPromise = $interval($scope.refreshJobs, 5000);
    $scope.$on('$destroy', function() {
        $interval.cancel(refreshPromise);
    });
});