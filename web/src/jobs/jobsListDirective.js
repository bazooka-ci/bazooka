"use strict";

angular.module('bzk.jobs').directive('bzkJobsList', function() {
  return {
    restrict: 'E',
    templateUrl: 'jobs/jobsList.html'
  };
});