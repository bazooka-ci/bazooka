angular.module('bzk.home').controller('HomeController', function($scope, HomeJobResource, $interval){
	$scope.refreshJobs = function() {
		HomeJobResource.jobs().success(function(jobs){
			$scope.jobs = jobs;
		});
	};

	$scope.refreshJobs();

	var refreshPromise = $interval($scope.refreshJobs, 5000);
	$scope.$on('$destroy', function() {
		$interval.cancel(refreshPromise);
	});
});
