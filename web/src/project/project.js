"use strict";

angular.module('bzk.project', ['bzk.jobs', 'bzk.utils', 'ngRoute']);

angular.module('bzk.project').config(function($routeProvider){
	$routeProvider.when('/p/:pid', {
			templateUrl: 'project/project.html',
			controller: 'ProjectController',
			reloadOnSearch: false
		});
});

angular.module('bzk.project').factory('ProjectResource', function($http){
	return {
		fetch: function(id) {
			return $http.get('/api/project/'+id);
		},
		jobs: function (id) {
			return $http.get('/api/project/'+id+'/job');
		},
		build: function (id, reference) {
			return $http.post('/api/project/'+id+'/job', {
				reference: reference
			});
		}
	};
});
