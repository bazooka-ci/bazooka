"use strict";

angular.module('bzk.projects', ['bzk.utils', 'ngRoute']);

angular.module('bzk.projects').config(function($routeProvider){
	$routeProvider.when('/p', {
			templateUrl: 'projects/projects.html',
			controller: 'ProjectListController',
			reloadOnSearch: false
		});
});

angular.module('bzk.projects').factory('ProjectListResource', function($http){
	return {
		list: function() {
			return $http.get('/api/project');
		},
		create: function (project) {
			return $http.post('/api/project', project);
		}
	};
});

angular.module('bzk.projects').controller('ProjectListController', function($scope, $routeParams, ProjectListResource){
	var pId = $routeParams.pid;

	ProjectListResource.list().success(function(projectList){
		$scope.projectList = projectList;
	});
});

angular.module('bzk.projects').controller('NewProjectController', function($scope, $routeParams, ProjectListResource){
	$scope.createProject = function(project) {
		ProjectListResource.create(project).success(function(jobs){
			$scope.jobs = jobs;
			console.log(jobs);
		});
	};
});
