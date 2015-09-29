"use strict";

angular.module('bzk', ['bzk.home', 'bzk.projects', 'bzk.project', 'bzk.job', 'bzk.utils', 'bzk.templates', 'ngRoute', 'angular-growl']);

angular.module('bzk').config(function($routeProvider) {
    $routeProvider.otherwise({
        redirectTo: '/'
    });
});

angular.module('bzk').config(function(growlProvider) {
    growlProvider.globalTimeToLive(3000);
    growlProvider.globalDisableCountDown(true);
});

angular.module('bzk').controller('RootController', function($scope, BzkApi, $routeParams, $location) {

});

angular.module('bzk').controller('ProjectsController', function($scope, BzkApi, EventBus, $routeParams) {

    function refresh() {
        BzkApi.project.list(true).success(function(res) {
            $scope.projects = res;
        });
    }


    EventBus.on('project.new', function() {
        refresh();
    });

    EventBus.on('jobs.refreshed', function(event, jobs) {
        var lastJobByProject = _(jobs).
        groupBy('project_id').
        mapValues(function(jobs) {
            return _(jobs).filter(function(job) {
                return job.status !== 'RUNNING';
            }).reduce(function(latest, job) {
                if (latest && job.number > latest.number) {
                    return job;
                }
                return latest;
            });
        }).value();

        _.forEach($scope.projects, function(project) {
            var job = lastJobByProject[project.id];
            if (job) {
                if (!project.last_job || job.number >= project.last_job.number) {
                    project.last_job = job;
                }
            }
        });
    });

    refresh();

    $scope.isSelected = function(p) {
        return p.id.indexOf($routeParams.pid) === 0;
    };

    $scope.projectClass = function(project) {
        return {
            'SUCCESS': true
        };
    };
});



angular.module('bzk').filter('bzoffset', function() {
    return function(o, b) {
        var t = o + b;
        if (t < 60) {
            return t + ' secs';
        } else if (t < 3600) {
            return Math.floor(t / 60) + ' mins';
        } else {
            return Math.floor(t / 3600) + ' hours';
        }
    };
});