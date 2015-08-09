'use strict';

angular.module('dashboard').controller('GridController', [
  '$scope',
  'widgetsManager',
  'LoadWidget',
  function($scope, widgetsManager, LoadWidget) {
    widgetsManager.register(new LoadWidget("webserver", 0, 0));

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
  }
]);

angular.module('dashboard').controller('WidgetController', [
  function() {
  }
]);
