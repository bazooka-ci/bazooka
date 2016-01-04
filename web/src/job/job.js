"use strict";

angular.module('bzk.job', ['bzk.utils', 'ngRoute']);

angular.module('bzk.job').config(function($routeProvider) {
    $routeProvider.when('/p/:pid/:jid/:vid?', {
        templateUrl: 'job/job.html',
        controller: 'JobController',
        reloadOnSearch: false
    });
});

angular.module('bzk.job').controller('JobLogsController', function($scope, BzkApi, DateUtils, $routeParams, $timeout, $interval) {
    var jId = $routeParams.jid,
        vId = $routeParams.vid;
    $scope.loadLogs = function(onNode, onDone) {
        var stream = vId ? BzkApi.variant.streamLog(vId, onNode, onDone) : BzkApi.job.streamLog(jId, onNode, onDone);
        $scope.$on('$destroy', stream.abort);
    };

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