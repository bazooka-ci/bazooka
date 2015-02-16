"use strict";

angular.module('bzk', ['bzk.projects', 'bzk.project', 'bzk.job', 'bzk.variant', 'bzk.utils', 'bzk.templates', 'ngRoute']);

angular.module('bzk').config(function($routeProvider){
	$routeProvider
	.when('/', {
		templateUrl: 'home.html',
		controller: 'HomeController',
	}).otherwise({
		redirectTo: '/'
	});
});

angular.module('bzk').factory('HomeJobResource', function($http){
	return {
		jobs: function() {
			return $http.get('/api/job');
		}
	};
});

angular.module('bzk').controller('HomeController', function($scope, HomeJobResource, $interval){
	$scope.refreshJobs = function() {
		HomeJobResource.jobs().success(function(jobs){
			$scope.jobs = jobs;
			console.log(jobs);
		});
	};

	$scope.refreshJobs();

	var refreshPromise = $interval($scope.refreshJobs, 5000);
	$scope.$on('$destroy', function() {
		$interval.cancel(refreshPromise);
	});
});

angular.module('bzk').controller('RootController', function($scope, $routeParams, $location, ProjectsResource){
	$scope.isProjectSelected = function() {
		return $routeParams.pid;
	};

	$scope.isJobSelected = function() {
		return $location.search().j;
	};

	$scope.isVariantSelected = function() {
		return $location.search().v;
	};

	$scope.getProjectName = function(projectID) {
		var project = _.findWhere($scope.all_projects, {id: projectID});
		if(project) {
			return project.name;
		}
		return projectID;
	};

	function refresh () {
		ProjectsResource.fetch().success(function(res){
			$scope.all_projects = res;
		});
	}

	refresh();
});

angular.module('bzk').factory('ProjectsResource', function($http){
	return {
		fetch: function () {
			return $http.get('/api/project');
		},
		create: function(proj) {
			return $http.post('/api/project', proj);
		}
	};
});

angular.module('bzk').controller('ProjectsController', function($scope, ProjectsResource, $routeParams){

	function refresh () {
		ProjectsResource.fetch().success(function(res){
			$scope.projects = res;
		});
	}

	refresh();

	$scope.newProj = {
		scm_type: 'git'
	};

	$scope.newProjectVisible = function(s) {
		$scope.showNewProject = s;
	};

	$scope.createProject = function() {
		ProjectsResource.create($scope.newProj).success(function(){
			$scope.showNewProject = false;
			$scope.newProj = {
				scm_type: 'git'
			};
			refresh();
		});
	};

	$scope.isSelected = function(p) {
		return p.id.indexOf($routeParams.pid)===0;
	};
});



angular.module('bzk').filter('bzoffset', function(){
	return function(o, b) {
		var t = o+b;
		if (t<60) {
			return t+' secs';
		} else if (t<3600) {
			return Math.floor(t/60)+' mins';
		} else {
			return Math.floor(t/3600) + ' hours';
		}
	};
});
