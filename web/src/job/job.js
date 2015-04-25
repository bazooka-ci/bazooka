"use strict";

angular.module('bzk.job', ['bzk.utils', 'ngRoute']);

angular.module('bzk.job').config(function($routeProvider) {
    $routeProvider.when('/p/:pid/:jid', {
        templateUrl: 'job/job.html',
        controller: 'JobController',
        reloadOnSearch: false
    });
});

angular.module('bzk.job').controller('JobLogsController', function($scope, BzkApi, DateUtils, $routeParams, $timeout) {
    var jId = $routeParams.jid;
    $scope.logger = {};

    function loadLogs() {
        $scope.logger.job.prepare();

        BzkApi.job.log(jId).success(function(logs) {
            $scope.logger.job.finish(logs);
        });
    }

    $timeout(loadLogs);
});