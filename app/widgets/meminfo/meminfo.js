'use strict';

angular.module('dashboard').controller('MeminfoWidgetController', [
  '$scope',
  function($scope) {
    $scope.isCritical = function() {
      return $scope.data.used > $scope.data.total / 2;
    };
  }
]);

angular.module('dashboard').factory('MeminfoWidget', [
  'Widget',
  function(Widget) {
    return function() {
      var self = new Widget('meminfo-widget', 'Memory');

      self.data = {
        total: 0,
        used: 0,
      };

      return self;
    };
  }
]);

angular.module('dashboard').directive('meminfoWidget', [
  function() {
    return {
      replace: true,
      templateUrl: 'widgets/meminfo/meminfo.html',
    };
  }
]);
