"use strict";

angular.module('bzk.home', ['bzk.utils', 'bzk.jobs', 'ngRoute']);

angular.module('bzk.home').config(function($routeProvider){
	$routeProvider.when('/', {
			templateUrl: 'home/home.html',
			controller: 'HomeController',
			reloadOnSearch: false
		});
});

angular.module('bzk.home').factory('HomeJobResource', function($http){
	return {
		jobs: function() {
			return $http.get('/api/job');
		}
	};
});
