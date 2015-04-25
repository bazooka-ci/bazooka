"use strict";

angular.module('bzk.utils').directive('bzkJobInfo', function() {
  return {
    restrict: 'AE',
    replace: true,
    scope: {
    	job: '&'
    },
    templateUrl: 'utils/job-info.html'
  };
});