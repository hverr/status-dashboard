'use strict';

angular.module('dashboard').controller('GridController', [
  '$scope',
  'widgetsManager',
  'Widget',
  function($scope, widgetsManager, Widget) {
    widgetsManager.register(new Widget('webserver', 'load', {
      sizeX: 1, sizeY: 1, row: 0, col: 0,
    }));

    $scope.gridsterOpts = {
      columns: 4,
    };

    $scope.widgets = [
      { sizeX: 1, sizeY: 1, row: 0, col: 0, data: { client : "Webserver", name : "Load", directive: "load-widget" }},
      { sizeX: 1, sizeY: 1, row: 0, col: 1, data: { client : "Webserver", name : "Uptime", directive: "uptime-widget" }},
    ];
  }
]);

angular.module('dashboard').controller('WidgetController', [
  '$log',
  '$scope',
  function($log, $scope) {
    $log.debug("WidgetController:", $scope.widget.data);
    $scope.widget = $scope.widget.data;
  }
]);
