"use strict";

angular.module('bzk.utils').directive('bzkJobInfo', function(BzkProjectsCache) {
    return {
        restrict: 'AE',
        scope: {
            job: '&bzkJobInfo',
            detailed: '&bzkDetailed'
        },
        controller: function($scope) {
            $scope.projectName = function(id) {
                var project = BzkProjectsCache.byId(id);
                if (project) {
                    return project.name;
                }
                return id;
            };
        },
        templateUrl: 'utils/job-info.html'
    };
});