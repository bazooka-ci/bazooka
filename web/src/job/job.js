"use strict";

angular.module('bzk.job', ['bzk.utils', 'ngRoute']);

angular.module('bzk.job').config(function($routeProvider) {
    $routeProvider.when('/p/:pid/:jid', {
        templateUrl: 'job/job.html',
        controller: 'JobController',
        reloadOnSearch: false
    });
});

angular.module('bzk.job').controller('JobLogsController', function($scope, BzkApi, DateUtils, $routeParams, $timeout, $interval) {
    var jId = $routeParams.jid;
    $scope.logger = {};

    function loadLogs() {
        $scope.logger.job.prepare();
        var stream = BzkApi.job.streamLog(
            jId,
            function(logEntry) {
                $scope.logger.job.append([logEntry]);
            },
            function() {
                $scope.logger.job.finish([]);
            }
        );

        $scope.$on('$destroy', stream.abort);
    }

    $timeout(loadLogs);
});