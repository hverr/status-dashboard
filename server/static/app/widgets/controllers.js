'use strict';

angular.module('dashboard').controller('GridController', [
  '$scope',
  'widgetsManager',
  'LoadWidget',
  'UptimeWidget',
  '$log',
  function($scope, widgetsManager, LoadWidget, UptimeWidget, $log) {
    $scope.addColumn = function() {
      $scope.gridsterOpts.columns += 1;
    };

    $scope.removeColumn = function() {
      if($scope.gridsterOpts.columns > 0) {
        $scope.gridsterOpts.columns -= 1;
      }
    };

    $scope.addWidgets = function() {
      $log.debug('addWidgets');
    };

    $scope.saveLayout = function() {
      $log.debug('saveLayout');
    };

    $scope.widgetGridsterMap = {
      sizeX: 'widget.height',
      sizeY: 'widget.width',
      row: 'widget.row',
      col: 'widget.col',
    };

    $scope.gridsterOpts = {
      columns: 4,
      margins: [16, 16],
    };

    widgetsManager.start();
  }
]);

angular.module('dashboard').controller('WidgetController', [
  function() {
  }
]);
