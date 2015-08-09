'use strict';

angular.module('dashboard').controller('LoadWidgetController', [
  function() {}
]);

angular.module('dashboard').factory('LoadWidget', [
  function() {
    return function(client, row, col) {
      return {
        directive: "load-widget",
        height: 1,
        width: 1,
        row: row,
        col: col,

        client: client,
        identifier : 'load',
        name: "Load",

        cores: 4,
        one : "1.05",
        five : "4.02",
        fifteen : "1.02",

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
