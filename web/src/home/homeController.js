"use strict";

angular.module('bzk.home').controller('HomeController', function($scope, BzkApi, EventBus, $interval) {
    $scope.refreshJobs = function() {
        BzkApi.job.list().success(function(jobs) {
            $scope.jobs = jobs;
            EventBus.send('jobs.refreshed', jobs);
        });
    };

    $scope.refreshJobs();

    var refreshPromise = $interval($scope.refreshJobs, 5000);
    $scope.$on('$destroy', function() {
        $interval.cancel(refreshPromise);
    });
});