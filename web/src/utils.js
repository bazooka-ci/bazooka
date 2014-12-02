"use strict";

angular.module('bzk.utils', []);

angular.module('bzk.utils').filter('bzkStatus', function(){
	var st2class = {
		'RUNNING': 'running',
		'SUCCESS': 'success',
		'FAILED': 'failed',
		'ERRORERD': 'errored'
	};

	return function(job) {
		return st2class[job.status];
	};

});

angular.module('bzk.utils').filter('bzkFinished', function(){
	return function(job) {
		var m = moment(job.completed);
		return m.year()==1? '-':m.format('HH:mm:ss - DD MMM YYYY');
	};
});

angular.module('bzk.utils').filter('bzkDate', function(){
	return function(d) {
		var m = moment(d);
		return m.year()==1? '-':m.format('HH:mm:ss - DD MMM YYYY');
	};
});


angular.module('bzk.utils').filter('bzkDuration', function(){
	return function(job) {
		var m = moment(job.completed);
		if (m.year()===1) {
			return moment().from(moment(job.started), true);
		} else {
			return m.from(job.started, true);
		}
	};
});

angular.module('bzk.utils').factory('bzkScroll', function($window){
	return {
		toTheRight: function(){
			$('html, body').animate({
				scrollLeft: 1000
			}, 500);
		}
	};
});
