"use strict";

angular.module('bzk.utils').factory('BzkApi', function($http, JsonStream) {
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
            build: function(id, reference, parameters) {
                return $http.post('/api/project/' + id + '/job', {
                    reference: reference,
                    parameters: parameters
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
            },
            streamLog: function(jid, onNode, onDone) {
                return JsonStream({
                    url: '/api/job/' + jid + '/log?follow=true&strict-json=true',
                    pattern: '{id}',
                    onNode: onNode,
                    onDone: onDone
                });
            }
        },
        variant: {
            log: function(vid) {
                return $http.get('/api/variant/' + vid + '/log');
            },
            streamLog: function(vid, onNode, onDone) {
                return JsonStream({
                    url: '/api/variant/' + vid + '/log?follow=true&strict-json=true',
                    pattern: '{id}',
                    onNode: onNode,
                    onDone: onDone
                });
            }
        }
    };
});