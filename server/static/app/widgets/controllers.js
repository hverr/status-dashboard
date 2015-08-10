'use strict';

angular.module('dashboard').controller('GridController', [
  '$scope',
  'widgetsManager',
  'LoadWidget',
  'UptimeWidget',
  function($scope, widgetsManager, LoadWidget, UptimeWidget) {
    widgetsManager.register(new LoadWidget("webserver", 0, 0));
    widgetsManager.register(new UptimeWidget("webserver", 0, 1));

    $scope.widgetGridsterMap = {
      sizeX: 'widget.height',
      sizeY: 'widget.width',
      row: 'widget.row',
      col: 'widget.col',
    };

    $scope.gridsterOpts = {
      columns: 4,
    };

    $scope.widgets = widgetsManager.registeredWidgets();

    widgetsManager.start();
  }
]);

angular.module('dashboard').controller('WidgetController', [
  function() {
  }
]);
