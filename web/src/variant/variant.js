"use strict";

angular.module('bzk.variant', ['bzk.utils', 'ngRoute']);

angular.module('bzk.variant').config(function($routeProvider) {
    $routeProvider.when('/p/:pid/:jid/:vid', {
        templateUrl: 'variant/variant.html',
        controller: 'VariantController',
        reloadOnSearch: false
    });
});

angular.module('bzk.variant').controller('VariantController', function($scope, BzkApi, EventBus, DateUtils, $routeParams, $timeout) {
    var pId = $routeParams.pid,
        jId = $routeParams.jid,
        vId = $routeParams.vid;

    var refreshPromise;

    $scope.$on('$destroy', function() {
        $timeout.cancel(refreshPromise);
    });

    function refresh() {
        BzkApi.project.get(pId).success(function(project) {
            $scope.project = project;
        });

        BzkApi.job.get(jId).success(function(job) {
            $scope.job = job;
            EventBus.send('jobs.refreshed', [job]);

            if (job.status === 'RUNNING') {
                refreshPromise = $timeout(refresh, 3000);
            }
        });

        BzkApi.variant.get(vId).success(function(variant) {
            $scope.variant = variant;
        });
    }
    refresh();
});

angular.module('bzk.variant').controller('VariantLogsController', function($scope, BzkApi, $routeParams, $timeout) {
    var vId = $routeParams.vid;
    $scope.logger = {};

    function loadLogs() {
        $scope.logger.variant.prepare();

        var stream = BzkApi.variant.streamLog(
            vId,
            function(logEntry) {
                $scope.logger.variant.append([logEntry]);
            },
            function() {
                $scope.logger.variant.finish([]);
            }
        );

        $scope.$on('$destroy', function() {
            stream.abort();
        });
    }

    $timeout(loadLogs);
});