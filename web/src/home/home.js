"use strict";

angular.module('bzk.home', ['bzk.utils', 'ngRoute']);

angular.module('bzk.home').config(function($routeProvider){
	$routeProvider.when('/', {
			templateUrl: 'home/home.html',
			controller: 'HomeController',
			reloadOnSearch: false
		});
});