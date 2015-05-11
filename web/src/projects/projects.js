"use strict";

angular.module('bzk.projects', ['bzk.utils', 'ngRoute']);

angular.module('bzk.projects').config(function($routeProvider) {
    $routeProvider.when('/p', {
        templateUrl: 'projects/projects.html',
        controller: 'ProjectListController',
        reloadOnSearch: false
    });
});

angular.module('bzk.projects').controller('ProjectListController', function($scope, BzkApi, $rootScope, $routeParams, growl) {
    var pId = $routeParams.pid;
    BzkApi.project.list().success(function(projectList) {
        $scope.projectList = projectList;
    });

    $scope.createProject = function() {
        BzkApi.project.create($scope.project).success(function() {
            growl.success('New project <strong>'+$scope.project.name+'</strong> created');
            $scope.project = {};

            BzkApi.project.list().success(function(projectList) {
                $scope.projectList = projectList;
                $rootScope.$broadcast('project.new');
            });
        });
    };
});

angular.module('bzk.projects').controller('NewProjectController', function($scope, $routeParams, BzkApi) {

});