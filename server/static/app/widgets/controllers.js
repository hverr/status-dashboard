'use strict';

angular.module('dashboard').controller('GridController', [
  '$scope',
  'widgetsManager',
  'LoadWidget',
  function($scope, widgetsManager, LoadWidget) {
    widgetsManager.register(new LoadWidget("webserver", 0, 0));

    $scope.widgetGridsterMap = {
      sizeX: 'item.height',
      sizeY: 'item.width',
      row: 'item.row',
      col: 'item.col',
    };

    $scope.gridsterOpts = {
      columns: 4,
    };

    $scope.widgets = widgetsManager.registeredWidgets();
  }
]);

angular.module('dashboard').controller('WidgetController', [
  function() {
  }
]);
