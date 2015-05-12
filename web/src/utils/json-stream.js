"use strict";

angular.module('bzk.utils').factory('JsonStream', function($q) {
    return function(params) {
        var stream = oboe({
            url: params.url
        }).node(params.pattern, function(node) {
            params.onNode(node);
            return oboe.drop;
        }).done(function() {
            params.onDone();
            return oboe.drop;
        });

        return {
            abort: function() {
                stream.abort();
            }
        };
    };
});