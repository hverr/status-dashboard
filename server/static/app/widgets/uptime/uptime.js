'use strict';

angular.module('dashboard').controller('UptimeWidgetController', [
  function() {
  }
]);

angular.module('dashboard').factory('UptimeWidget', [
  'Widget',
  function(Widget) {
    return function() {
      var self = new Widget('uptime-widget', 'Uptime');

      self.data = {
        days: 0,
        hours: 0,
        minutes: 0,
        seconds: 0,
      };

      // We clone the data, because otherwise multiple directives sharing
      // the data would increase the time more than once per second.
      self.update = function(data) {
        self.data.days = data.days;
        self.data.hours = data.hours;
        self.data.minutes = data.minutes;
        self.data.seconds = data.seconds;
      };

      return self;
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
        if(!scope.data) {
          return;
        }

        scope.data.seconds += 1;

        if(scope.data.seconds >= 60) {
          scope.data.seconds = 0;
          scope.data.minutes += 1;
        }

        if(scope.data.minutes >= 60) {
          scope.data.minutes = 0;
          scope.data.hours += 1;
        }

        if(scope.data.hours >= 24) {
          scope.data.hours = 0;
          scope.data.days += 1;
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
