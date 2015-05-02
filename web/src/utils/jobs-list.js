"use strict";

angular.module('bzk.utils').directive('bzkJobsList', function() {
    return {
        restrict: 'AE',
        scope: {
            jobs: '&'
        },
        templateUrl: 'utils/jobs-list.html',
        controller: function($scope, BzkApi) {
            $scope.only = function(status) {
                return function(job) {
                    return job.status === status;
                };
            };

            $scope.not = function(status) {
                return function(job) {
                    return job.status !== status;
                };
            };

            BzkApi.project.list().success(function(projects) {
                $scope.projectName = function(projectID) {
                    var project = _.findWhere(projects, {
                        id: projectID
                    });
                    if (project) {
                        return project.name;
                    }
                    return projectID;
                };
            });
        }
    };
});