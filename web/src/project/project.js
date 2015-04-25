"use strict";

angular.module('bzk.project', ['bzk.utils', 'ngRoute']);

angular.module('bzk.project').config(function($routeProvider) {
    $routeProvider.when('/p/:pid', {
        templateUrl: 'project/project.html',
        controller: 'ProjectController',
        reloadOnSearch: false
    });
});