"use strict";

angular.module('bzk.utils', []);

angular.module('bzk.utils').filter('bzkStatus', function() {
    var st2class = {
        'RUNNING': 'running',
        'SUCCESS': 'success',
        'FAILED': 'failed',
        'ERRORED': 'errored'
    };

    return function(job) {
        return st2class[job.status];
    };

});

angular.module('bzk.utils').factory('DateUtils', function() {
    return {
        isSet: function(date) {
            var m = moment(date);
            return m.year() != 1;
        }
    };
});

angular.module('bzk.utils').filter('bzkFinished', function() {
    return function(job) {
        if (job) {
            var m = moment(job.completed);
            return m.year() == 1 ? '-' : m.format('HH:mm:ss - DD MMM YYYY');
        }
    };
});

angular.module('bzk.utils').filter('bzkDate', function() {
    return function(d) {
        if (d) {
            var m = moment(d);
            return m.year() == 1 ? '-' : m.format('HH:mm:ss - DD MMM YYYY');
        }
    };
});

angular.module('bzk.utils').filter('bzkDuration', function() {
    return function(job) {
        if (job) {
            var m = moment(job.completed);
            if (m.year() === 1) {
                return moment().from(moment(job.started), true);
            } else {
                return m.from(job.started, true);
            }
        }
    };
});

angular.module('bzk.utils').filter('bzkExcerpt', function() {
    return function(s, width) {
        width = width || 7;
        if (s) {
            return s.substr(0, width);
        }
    };
});

angular.module('bzk.utils').filter('bzkId', function($filter) {
    return function(obj) {
        if (obj) {
            return $filter('bzkExcerpt')(obj.id, 7);
        }
    };
});

angular.module('bzk.utils').filter('titleCase', function() {
    return function(str) {
        return (str === undefined || str === null) ? '' : str.replace(/_|-/, ' ').replace(/\w\S*/g, function(txt) {
            return txt.charAt(0).toUpperCase() + txt.substr(1).toLowerCase();
        });
    };
});

angular.module('bzk.utils').factory('Scroll', function($window) {
    return {
        toTheRight: function() {
            $('html, body').animate({
                scrollLeft: 1000
            }, 500);
        }
    };
});

angular.module('bzk.utils').directive('bzkLog', function() {
    return {
        restrict: 'A',
        scope: {
            sink: '=bzkLog'
        },
        template: '<img class="loading" src="/images/loading.gif" ng-if="loading"></img>',
        link: function($scope, elem, attrs) {
            var row = 1;
            $(elem).append('<pre></pre>');
            var into = $(elem).find('pre');
            $scope.sink = {
                prepare: function() {
                    this.clear();
                    $scope.loading = true;
                },
                finish: function(lines) {
                    this.append(lines);
                    $scope.loading = false;
                },
                append: function(lines) {
                    var data = '';
                    _.each(lines, function(line) {
                        data += '<p><span>' + row + '</span>' + _.escape(line.msg) + '</p>';
                        row++;
                    });
                    into.append(data);
                },
                clear: function() {
                    row = 1;
                    into.empty();
                    into.scrollTop(0);
                }
            };
        }
    };
});