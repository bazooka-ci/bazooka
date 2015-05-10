"use strict";

angular.module('bzk.variant', ['bzk.utils', 'ngRoute']);

angular.module('bzk.variant').config(function($routeProvider) {
    $routeProvider.when('/p/:pid/:jid/:vid', {
        templateUrl: 'variant/variant.html',
        controller: 'VariantController',
        reloadOnSearch: false
    });
});

angular.module('bzk.variant').controller('VariantController', function($scope, BzkApi, DateUtils, $routeParams, $timeout) {
    var jId;
    var pId;
    var vId;
    var refreshPromise;

    $scope.$on('$destroy', function() {
        $timeout.cancel(refreshPromise);
    });

    function refresh() {
        pId = $routeParams.pid;
        if (pId) {
            BzkApi.project.get(pId).success(function(project) {
                $scope.project = project;
            });
        }
        jId = $routeParams.jid;
        if (jId) {
            BzkApi.job.get(jId).success(function(job) {
                $scope.job = job;

                if (job.status === 'RUNNING') {
                    refreshPromise = $timeout(refresh, 3000);
                }
            });

            BzkApi.job.variants(jId).success(function(variants) {
                var result = $.grep(variants, function(e) {
                    return e.id.indexOf($routeParams.vid) === 0;
                });
                if (result) {
                    $scope.variant = result[0];
                }
            });
        }
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