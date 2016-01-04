"use strict";

angular.module('bzk.utils', []);

angular.module('bzk.utils').factory('EventBus', function($rootScope) {
    return {
        send: function(key, value) {
            $rootScope.$emit(key, value);
        },
        on: function(key, callback) {
            return $rootScope.$on(key, function(event, args) {
                callback(event.name, args);
            });
        }
    };
});

angular.module('bzk.utils').factory('BzkProjectsCache', function(BzkApi) {
    var projects;

    var cache = {
        byId: function(projectID) {
            return _.findWhere(projects, {
                id: projectID
            });
        },
        refresh: function() {
            BzkApi.project.list().success(function(res) {
                projects = res;
            });
        }
    };

    cache.refresh();
    return cache;
});


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
            return m.year() == 1 ? '-' : m.format('HH:mm - DD MMM YYYY');
        }
    };
});

angular.module('bzk.utils').filter('bzkDuration', function() {
    function fmt(from, to) {
        return moment.duration(from.diff(to), 'milliseconds').format('h [hours] m [mins] s [secs]');
    }

    return function(job) {
        if (job) {
            var m = moment(job.completed);
            if (m.year() === 1) {
                return fmt(moment(), job.started);
            } else {
                return fmt(m, job.started);
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

angular.module('bzk.utils').directive('bzkLog', function($interval) {
    return {
        restrict: 'A',
        scope: {
            loader: '&bzkLog'
        },
        template: '<pre></pre>',
        link: function($scope, elem, attrs) {

            function isAtBottom() {
                return false; //$(window).scrollTop() + $(window).height() == $(document).height();
            }

            var row = 1,
                into = $(elem).find('pre'),
                marker;

            var periodicFlusher;
            function prepare() {
                clear();
                into.append('<p class="loading-marker"><span></span><img src="/images/loading.gif"></img></p>');
                marker = $(elem).find('.loading-marker');
                periodicFlusher = $interval(flush, 500);
            }

            function finish() {
                flush();
                $interval.cancel(periodicFlusher);
                marker.remove();
            }

            function append(line) {
                var scroll = isAtBottom(),
                    data = '';

                data += '<p><span>' + row + '</span>' + (_.escape(line.msg) || '&nbsp;') + '</p>';
                row++;

                $(data).insertBefore(marker);

                if (scroll) {
                    $('body').scrollTop(into.height());
                }
            }

            var lines = [];
            function lazy_append(line) {
                lines.push(line);
            }

            function flush() {
                var scroll = isAtBottom(),
                    data = '';

                _.each(lines, function(line) {
                    data += '<p><span>' + row + '</span>' + (_.escape(line.msg) || '&nbsp;') + '</p>';
                    row++;
                });

                $(data).insertBefore(marker);
                lines = [];

                if (scroll) {
                    $('body').scrollTop(into.height());
                }
            }

            function clear() {
                row = 1;
                into.empty();
                into.scrollTop(0);
            }


            prepare();
            $scope.loader()(lazy_append, finish);

        }
    };
});