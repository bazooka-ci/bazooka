"use strict";

angular.module('bzk.utils').directive('bzkBreadcrumbs', function(){
	return {
		restrict: 'AE',
		replace: true,
		scope: {
			project: '&',
			job: '&',
			variant: '&'
		},
		templateUrl: 'utils/breadcrumbs.html'
	};
});
