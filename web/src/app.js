"use strict";

angular.module('bzk', ['bzk.home', 'bzk.projects', 'bzk.project', 'bzk.job', 'bzk.variant', 'bzk.utils', 'bzk.templates', 'ngRoute']);

angular.module('bzk').config(function($routeProvider) {
    $routeProvider.otherwise({
        redirectTo: '/'
    });
});

angular.module('bzk').controller('RootController', function($scope, BzkApi, $routeParams, $location) {

});

angular.module('bzk').controller('ProjectsController', function($scope, BzkApi, $routeParams) {

    function refresh() {
        BzkApi.project.list().success(function(res) {
            $scope.projects = res;
        });
    }

    $scope.$on('project.new', function(event) {
        refresh();
    });

    refresh();

    $scope.isSelected = function(p) {
        return p.id.indexOf($routeParams.pid) === 0;
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