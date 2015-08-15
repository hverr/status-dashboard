'use strict';

angular.module('dashboard').controller('LoadWidgetController', [
  '$scope',
  function($scope) {
    $scope.isLargeLoad = function(load) {
      var f = parseFloat(load);
      if(isNaN(f)) {
        return false;
      }

      return f > $scope.widget.cores;
    };
  }
]);

angular.module('dashboard').factory('LoadWidget', [
  function() {
    return function() {
      return {
        directive: "load-widget",
        height: 1,
        width: 1,
        row: 0,
        col: 0,

        client: null,
        name: "Load",

        data : {
          cores: 4,
          one : "1.05",
          five : "4.02",
          fifteen : "1.02",
        },

        update : function(object) {
          this.cores = object.cores;
          this.one = object.one;
          this.five = object.five;
          this.fifteen = object.fifteen;
        },
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
