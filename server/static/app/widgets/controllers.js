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

    var l = widgetsManager.add('webserver', 'load');
    l.client = 'Web Server';
    $scope.widgets = [l];

    var u = widgetsManager.add('webserver', 'uptime');
    u.client = 'Web Server';
    $scope.widgets.push(u);

    widgetsManager.update(true);
  }
]);

angular.module('dashboard').controller('WidgetDataController', [
  '$scope',
  function($scope) {
    $scope.data = $scope.widget.data;
    $scope.$watch('widget.data', function(newValue) {
      $scope.data = newValue;
    });
  }
]);

angular.module('dashboard').controller('AddWidgetsDialogController', [
  '$scope',
  '$rootScope',
  'widgetsManager',
  '$log',
  function($scope, $rootScope, widgetsManager, $log) {
    function update() {
      $scope.clients = widgetsManager.availableClients;
      $log.debug('Available:', $scope.clients);
      if(!$scope.clients) {
        $scope.message = "No clients connected, please refersh.";
      } else {
        $scope.message = null;
      }
    }

    update();
    $rootScope.$on(widgetsManager.availableClientsChangedEvent, function() {
      update();
    });

    $scope.addWidget = function(clientIdentifier, widgetType) {
      $log.debug('Should add', clientIdentifier, ':', widgetType);
    };
  }
]);
