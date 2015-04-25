"use strict";

angular.module('bzk.utils').factory('BzkApi', function($http) {
    return {
        project: {
            list: function() {
                return $http.get('/api/project');
            },
            get: function(id) {
                return $http.get('/api/project/' + id);
            },
            create: function(project) {
                return $http.post('/api/project', project);
            },
            jobs: function(id) {
                return $http.get('/api/project/' + id + '/job');
            },
            build: function(id, reference) {
                return $http.post('/api/project/' + id + '/job', {
                    reference: reference
                });
            }
        },
        job: {
        	list: function() {
                return $http.get('/api/job');
            },
            get: function(id) {
                return $http.get('/api/job/' + id);
            },
            variants: function(jid) {
                return $http.get('/api/job/' + jid + '/variant');
            },
            log: function(jid) {
                return $http.get('/api/job/' + jid + '/log');
            }
        },
        variant: {
            log: function(vid) {
                return $http.get('/api/variant/' + vid + '/log');
            }
        }
    };
});