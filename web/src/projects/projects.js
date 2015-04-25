"use strict";

angular.module('bzk.projects', ['bzk.utils', 'ngRoute']);

angular.module('bzk.projects').config(function($routeProvider){
	$routeProvider.when('/p', {
			templateUrl: 'projects/projects.html',
			controller: 'ProjectListController',
			reloadOnSearch: false
		});
});

angular.module('bzk.projects').controller('ProjectListController', function($scope, BzkApi, $rootScope, $routeParams){
	var pId = $routeParams.pid;

	BzkApi.project.list().success(function(projectList){
		$scope.projectList = projectList;
	});

	$scope.createProject = function(project) {
		BzkApi.project.create(project).success(function(){
			BzkApi.project.list().success(function(projectList){
				$scope.projectList = projectList;
				$rootScope.$broadcast('project.new');
			});
		});
	};
});

angular.module('bzk.projects').controller('NewProjectController', function($scope, $routeParams, BzkApi){

});
