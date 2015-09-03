'use strict';

angular.module('dashboard').controller('CurrentDateWidgetController', [
  function() {
  }
]);

angular.module('dashboard').factory('CurrentDateWidget', [
  'Widget',
  function(Widget) {
    return function() {
      var self = new Widget('current-date-widget', 'Date');

      self.available = true;
      self.data = {
        now: new Date(),
      };

      return self;
    };
  }
]);

angular.module('dashboard').directive('currentDateWidget', [
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
      templateUrl: 'widgets/current_date/current_date.html',
      link: link,
    };
  }
]);
