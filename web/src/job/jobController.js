"use strict";

angular.module('bzk.job').controller('JobController', function($scope, BzkApi, EventBus, $routeParams, $timeout) {
    var refreshJobPromise, refreshVariantsPromise;

    $scope.$on('$destroy', function() {
        $timeout.cancel(refreshJobPromise);
        $timeout.cancel(refreshVariantsPromise);
    });

    $scope.loadLogs = function(onNode, onDone) {
        var jId = $routeParams.jid,
            vId = $routeParams.vid;
        var stream = vId ? BzkApi.variant.streamLog(vId, onNode, onDone) : BzkApi.job.streamLog(jId, onNode, onDone);
        $scope.$on('$destroy', stream.abort);
    };

    function refresh() {
        BzkApi.project.get($routeParams.pid).success(function(project) {
            $scope.project = project;
        });

        BzkApi.job.get($routeParams.jid).success(function(job) {
            if(!$scope.job) {
                // first time we finished loading a job, load the variants
                refreshVariants();
            }
            $scope.job = job;
            EventBus.send('jobs.refreshed', [job]);

            if (job.status === 'RUNNING') {
                refreshJobPromise = $timeout(refresh, 3000);
            }
        });
    }

    refresh();

    function refreshVariants() {
        var jId = $routeParams.jid;
        if (jId) {
            BzkApi.job.variants(jId).success(function(variants) {
                $scope.variants = variants;

                var vId = $routeParams.vid;
                if(vId) {
                    var selected =  new RegExp('^'+vId);
                    $scope.selectedVariant = _.find(variants, function(v){
                        return selected.test(v.id);
                    });
                }

                if ($scope.job.status === 'RUNNING' || _.findWhere($scope.variants, {
                    status: 'RUNNING'
                })) {
                    refreshVariantsPromise = $timeout(refreshVariants, 3000);
                }
            });
        }
    }
});