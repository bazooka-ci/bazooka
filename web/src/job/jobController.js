"use strict";

angular.module('bzk.job').controller('JobController', function($scope, BzkApi, $routeParams, $timeout) {
    var jId;
    var pId;
    var refreshJobPromise, refreshVariantsPromise;

    $scope.$on('$destroy', function() {
        $timeout.cancel(refreshJobPromise);
        $timeout.cancel(refreshVariantsPromise);
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
                    refreshJobPromise = $timeout(refresh, 3000);
                }
            });
        }
    }

    refresh();

    function refreshVariants() {
        var jId = $routeParams.jid;
        if (jId) {
            BzkApi.job.variants(jId).success(function(variants) {

                $scope.variants = variants;
                setupMeta(variants);

                if ($scope.job.status === 'RUNNING' || _.findWhere($scope.variants, {
                    status: 'RUNNING'
                })) {
                    refreshVariantsPromise = $timeout(refreshVariants, 3000);
                }
            });
        }
    }

    $scope.$watch('job', function(newValue, oldValue) {
        if (!oldValue) {
            // first time we finished loading a job, load the variants
            refreshVariants();
        }
    });

    $scope.variantsStatus = function() {
        if ($scope.variants && $scope.variants.length > 0) {
            return 'show';
        } else if ($scope.job.status === 'RUNNING') {
            return 'pending';
        } else {
            return 'none';
        }
    };

    function setupMeta(variants) {
        var colorsDb = ['#4a148c' /* Purple */ ,
            '#006064' /* Cyan */ ,
            '#f57f17' /* Yellow */ ,
            '#e65100' /* Orange */ ,
            '#263238' /* Blue Grey */ ,
            '#b71c1c' /* Red */ ,
            '#1a237e' /* Indigo */ ,
            '#1b5e20' /* Green */ ,
            '#33691e' /* Light Green */ ,
            '#212121' /* Grey 500 */ ,
            '#880e4f' /* Pink */ ,
            '#311b92' /* Deep Purple */ ,
            '#01579b' /* Light Blue */ ,
            '#004d40' /* Teal */ ,
            '#ff6f00' /* Amber */ ,
            '#bf360c' /* Deep Orange */ ,
            '#0d47a1' /* Blue */ ,
            '#827717' /* Lime */ ,
            '#3e2723' /* Brown 500 */ ,
            '#000000'
        ];

        var metaList = [],
            colors = {};
        if (variants.length > 0) {
            var visitedMeta = {};
            _.each(variants, function(v) {
                v.metaMap = _.indexBy(v.metas, 'name');
                _.each(v.metas, function(m) {
                    if (!visitedMeta[m.name]) {
                        metaList.push(m);
                        m.label = m.kind == 'env' ? '$' + m.name : m.name;
                        visitedMeta[m.name] = true;
                    }
                });
            });

            _.each(metaList, function(m) {
                var mcolors = {};
                colors[m.name] = mcolors;
                var colIdx = 0;
                _.each(variants, function(v) {
                    var vm = _.findWhere(v.metas, {
                        name: m.name
                    });
                    if (vm) {
                        var val = vm.value;
                        if (!mcolors[val]) {
                            mcolors[val] = colorsDb[colIdx];
                            if (colIdx < colorsDb.length - 1) {
                                colIdx++;
                            }
                        }
                    }
                });
            });
        }

        $scope.metaList = metaList;
        $scope.metaColors = colors;
    }

    $scope.metaColor = function(vmeta) {
        if (vmeta) {
            return $scope.metaColors[vmeta.name][vmeta.value];
        }
    };
});