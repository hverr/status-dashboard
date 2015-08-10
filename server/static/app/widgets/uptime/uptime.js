'use strict';

angular.module('dashboard').controller('UptimeWidgetController', [
  function() {
  }
]);

angular.module('dashboard').factory('UptimeWidget', [
  function() {
    return function(client, row, col) {
      return {
        directive: "uptime-widget",
        height: 1,
        width: 1,
        row: row,
        col: col,

        client: client,
        identifier : 'uptime',
        name : "Uptime",

        days : 0,
        hours : 0,
        minutes : 0,
        seconds : 0,

        update : function(object) {
          this.days = object.days;
          this.hours = object.hours;
          this.minutes = object.minutes;
          this.seconds = object.seconds;
        },
      };
    };
  }
]);

angular.module('dashboard').directive('uptimeWidget', [
  '$interval',
  function($interval) {
    var timer;

    function link(scope, element) {
      element.on('$destroy', function() {
        $interval.cancel(timer);
      });

      var increase = function() {
        scope.widget.seconds += 1;

        if(scope.widget.seconds >= 60) {
          scope.widget.seconds = 0;
          scope.widget.minutes += 1;
        }

        if(scope.widget.minutes >= 60) {
          scope.widget.minutes = 0;
          scope.widget.hours += 1;
        }

        if(scope.widget.hours >= 24) {
          scope.widget.hours = 0;
          scope.widget.days += 1;
        }
      };

      timer = $interval(increase, 1000);
    }

    return {
      replace: true,
      templateUrl : 'widgets/uptime/uptime.html',
      link: link,
    };
  }
]);
