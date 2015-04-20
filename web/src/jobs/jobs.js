"use strict";

angular.module('bzk.jobs', ['bzk.utils', 'ngRoute']);

angular.module('bzk.jobs').config(function($routeProvider){
	$routeProvider.when('/p/:pid/:jid', {
			templateUrl: 'jobs/detail.html',
			controller: 'JobDetailController',
			reloadOnSearch: false
		});
});

angular.module('bzk.jobs').factory('JobResource', function($http){
	return {
		project: function(id) {
			return $http.get('/api/project/'+id);
		},
		job: function (id) {
			return $http.get('/api/job/'+id);
		},
		variants: function (jid) {
			return $http.get('/api/job/'+jid+'/variant');
		},
		jobLog: function (jid) {
			return $http.get('/api/job/'+jid+'/log');
		},
		variantLog: function (vid) {
			return $http.get('/api/variant/'+vid+'/log');
		}
	};
});
