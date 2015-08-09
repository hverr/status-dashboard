'use strict';

angular.module('dashboard').controller('GridController', [
  '$scope',
  'widgetsManager',
  'Widget',
  function($scope, widgetsManager, Widget) {
    widgetsManager.register(new Widget('webserver', 'load', {
      sizeX: 1, sizeY: 1, row: 0, col: 0,
    }));

    $scope.widgets = [
      { sizeX: 1, sizeY: 1, row: 0, col: 0, data: "hay"},
      { sizeX: 1, sizeY: 1, row: 0, col: 1, data: "joo"},
    ];
  }
]);

angular.module('dashboard').controller('WidgetController', [
  '$log',
  '$scope',
  function($log, $scope) {
    $log.debug("WidgetController:", $scope.widget);
  }
]);
