'use strict';

angular.module('dashboard').controller('CurrentTimeWidgetController', [
  function() {
  }
]);

angular.module('dashboard').factory('CurrentTimeWidget', [
  'Widget',
  function(Widget) {
    return function() {
      var self = new Widget('current-time-widget', 'Time');

      self.available = true;
      self.data = {
        now: new Date(),
      };

      return self;
    };
  }
]);

angular.module('dashboard').directive('currentTimeWidget', [
  'oneSecondService',
  function(oneSecondService) {
    function link(scope, element) {
      var handle = oneSecondService.add(function() {
        scope.data.now = new Date();
      });

      element.on('$destroy', function() {
        oneSecondService.remove(handle);
      });
    }

    return {
      replace: true,
      templateUrl : 'widgets/current_time/current_time.html',
      link: link,
    };
  }
]);
