"use strict";

angular.module('bzk.utils').directive('bzkJobStatus', function() {
    return {
        restrict: 'AE',
        replace: true,
        scope: {
            status: '&bzkJobStatus'
        },
        templateUrl: 'utils/job-status.html',
        controller: function($scope){
        	$scope.glyph = {
        		'SUCCESS': 'ok-circle',
        		'FAILED': 'remove-circle',
        		'ERRORED': 'ban-circle',
                'RUNNING': 'time'
        	};
        }
    };
});