'use strict';

angular.module('dashboard').controller('LoadWidgetController', [
  '$scope',
  function($scope) {
    $scope.cores = "0";
    $scope.one = "3.4";
    $scope.five = "1.5";
    $scope.fifteen = "0.5";
  }
]);

angular.module('dashboard').factory('LoadWidget', [
  function() {
    return function() {
      return {
        one : null,
        five : null,
        fifteen : null,
      };
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
