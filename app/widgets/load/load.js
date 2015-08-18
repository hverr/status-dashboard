'use strict';

angular.module('dashboard').controller('LoadWidgetController', [
  '$scope',
  function($scope) {
    $scope.isLargeLoad = function(load) {
      var f = parseFloat(load);
      if(isNaN(f)) {
        return false;
      }

      return f > $scope.data.cores;
    };
  }
]);

angular.module('dashboard').factory('LoadWidget', [
  'Widget',
  function(Widget) {
    return function() {
      var self = new Widget('load-widget', 'Load');

      self.data = {
        cores: 4,
        one : "1.05",
        five : "4.02",
        fifteen : "1.02",
      };

      return self;
    };
  }
]);

angular.module('dashboard').directive('loadWidget', [
  function() {
    return {
      replace: true,
      templateUrl : 'widgets/load/load.html',
    };
  }
]);
