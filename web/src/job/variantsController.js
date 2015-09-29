"use strict";

angular.module('bzk.job').controller('VariantsController', function($scope, BzkApi, BzkColors) {
    $scope.$watch('variants', function(variants) {
        if (variants) {
            setupMeta(variants);
        }
    });

    $scope.isSelected = function(variant) {
        return $scope.selected().id === variant.id;
    };

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
                        m.len = m.label.length;
                        visitedMeta[m.name] = true;
                    }
                });
            });

            _.each(metaList, function(m) {
                var mcolors = {};
                colors[m.name] = mcolors;
                var colIdx = 0;
                var longestValue = 0;
                _.each(variants, function(v) {
                    var vm = _.findWhere(v.metas, {
                        name: m.name
                    });
                    if (vm) {
                        var val = vm.value;
                        longestValue = val.length > longestValue ? val.length : longestValue;
                        if (!mcolors[val]) {
                            mcolors[val] = BzkColors[colIdx];
                            if (colIdx < BzkColors.length - 1) {
                                colIdx++;
                            }
                        }
                    }
                });
                m.len = 30 + 9 * (m.len + longestValue);
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

angular.module('bzk.utils').value('BzkColors', ['#4a148c' /* Purple */ ,
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
]);