"use strict";

angular.module('bzk.variant', ['bzk.utils', 'ngRoute']);

angular.module('bzk.variant').config(function($routeProvider){
	$routeProvider.when('/p/:pid/:jid/:vid', {
			templateUrl: 'variant/variant.html',
			controller: 'VariantController',
			reloadOnSearch: false
		});
});

angular.module('bzk.variant').factory('VariantResource', function($http){
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

angular.module('bzk.variant').controller('VariantController', function($scope, VariantResource, DateUtils, $routeParams, $timeout){
	var jId;
	var pId;
	var vId;
	var refreshPromise;

	$scope.$on('$destroy', function() {
		$timeout.cancel(refreshPromise);
	});

	function refresh() {
		pId = $routeParams.pid;
		if(pId) {
			VariantResource.project(pId).success(function(project){
				$scope.project = project;
			});
		}
		jId = $routeParams.jid;
		if(jId) {
			VariantResource.job(jId).success(function(job){
				$scope.job = job;

				if (job.status==='RUNNING') {
					refreshPromise = $timeout(refresh, 3000);
				}
			});

			VariantResource.variants(jId).success(function(variants){
				var result = $.grep(variants, function(e){ return e.id.indexOf($routeParams.vid) === 0; });
				if(result) {
					$scope.variant = result[0];
				}
			});
		}
	}
	refresh();

	$scope.$on('$routeUpdate', refresh);
});

angular.module('bzk.variant').controller('VariantLogsController', function($scope, VariantResource, $routeParams, $timeout){
	var vId = $routeParams.vid;
	$scope.logger={};
	function loadLogs() {
		$scope.logger.variant.prepare();

		VariantResource.variantLog(vId).success(function(logs){
			$scope.logger.variant.finish(logs);
		});
	}

	$timeout(loadLogs);

});
