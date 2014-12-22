"use strict";

angular.module('bzk.job', ['bzk.utils', 'ngRoute']);

angular.module('bzk.job').config(function($routeProvider){
	$routeProvider.when('/p/:pid/:jid', {
			templateUrl: 'job/job.html',
			controller: 'JobController',
			reloadOnSearch: false
		});
});

angular.module('bzk.job').factory('JobResource', function($http){
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

angular.module('bzk.job').controller('JobController', function($scope, JobResource, DateUtils, $routeParams, $timeout){
	var jId;
	var pId;
	var refreshPromise;

	$scope.$on('$destroy', function() {
		$timeout.cancel(refreshPromise);
	});

	function refresh() {
		pId = $routeParams.pid;
		if(pId) {
			JobResource.project(pId).success(function(project){
				$scope.project = project;
			});
		}
		jId = $routeParams.jid;
  		if(jId) {
			JobResource.job(jId).success(function(job){
				$scope.job = job;

				if (job.status==='RUNNING') {
					refreshPromise = $timeout(refresh, 3000);
				}
			});
		}
	}

	refresh();

	$scope.$on('$routeUpdate', refresh);

});

angular.module('bzk.job').controller('JobLogsController', function($scope, JobResource, DateUtils, $routeParams, $timeout){
	var jId = $routeParams.jid;
	$scope.logger={};
	function loadLogs() {
		$scope.logger.job.prepare();

		JobResource.jobLog(jId).success(function(logs){
			$scope.logger.job.finish(logs);
		});
	}

	$timeout(loadLogs);

});

angular.module('bzk.job').controller('VariantsController', function($scope, JobResource, $routeParams, $timeout){
	var refreshPromise;

	$scope.$on('$destroy', function() {
		$timeout.cancel(refreshPromise);
	});

	function refreshVariants() {
		var jId = $routeParams.jid;
  		if(jId) {
			JobResource.variants(jId).success(function(variants){

				$scope.variants = variants;
				setupMeta(variants);

				if($scope.job.status==='RUNNING' || _.findWhere($scope.variants, {status: 'RUNNING'})) {
					refreshPromise= $timeout(refreshVariants, 3000);
				}
			});
		}
	}

	refreshVariants();

	$scope.variantsStatus = function() {
		if($scope.variants && $scope.variants.length>0) {
			return 'show';
		} else if($scope.job.status==='RUNNING'){
			return 'pending';
		} else {
			return 'none';
		}
	};

	$scope.$on('$routeUpdate', function(){
		refreshVariants();
	});

	function setupMeta(variants) {
		var colorsDb = ['#4a148c' /* Purple */,
	'#006064' /* Cyan */,
	'#f57f17' /* Yellow */,
	'#e65100' /* Orange */,
	'#263238' /* Blue Grey */,
	'#b71c1c' /* Red */,
	'#1a237e' /* Indigo */,
	'#1b5e20' /* Green */,
	'#33691e' /* Light Green */,
	'#212121' /* Grey 500 */,
	'#880e4f' /* Pink */,
	'#311b92' /* Deep Purple */,
	'#01579b' /* Light Blue */,
	'#004d40' /* Teal */,
	'#ff6f00' /* Amber */,
	'#bf360c' /* Deep Orange */,
	'#0d47a1' /* Blue */,
	'#827717' /* Lime */,
	'#3e2723' /* Brown 500 */,
	'#000000'];

		var metaLabels = [], colors={};
		if (variants.length>0) {
			var vref = variants[0];
			_.each(vref.metas, function (m) {
				metaLabels.push(m.kind=='env'?'$'+m.name:m.name);
			});

			_.each(vref.metas, function(m, i){
				var mcolors={};
				colors[m.name] = mcolors;
				var colIdx=0;
				_.each(variants, function (v) {
					var val=v.metas[i].value;
					if (!mcolors[val]) {
						mcolors[val] = colorsDb[colIdx];
						if(colIdx<colorsDb.length-1) {
							colIdx++;
						}
					}
				});
			});

		}

		$scope.metaLabels=metaLabels;
		$scope.metaColors=colors;
	}

	$scope.metaColor = function(vmeta) {
		return $scope.metaColors[vmeta.name][vmeta.value];
	};
});
