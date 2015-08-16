'use strict';

angular.module('dashboard').controller('GridController', [
  '$scope',
  '$rootScope',
  '$window',
  '$location',
  'widgetsManager',
  'LoadWidget',
  'UptimeWidget',
  '$log',
  function($scope, $rootScope, $window, $location, widgetsManager, LoadWidget, UptimeWidget, $log) {
    function findFreeTile(width) {
      var cols = $scope.gridsterOpts.columns;
      var row = 0;
      var col = 0;

      $scope.widgets.forEach(function(w) {
        if(w.row > row) {
          row = w.row;
          col = 0;
        }

        if(w.row === row && w.col + w.width > col) {
          col = w.col + w.width;
        }

        if(col >= cols) {
          row += 1;
          col = 0;
        }
      });

      if(col !== 0 && col + width > cols) {
        col = 0;
        row += 1;
      }

      return {
        row: row,
        column: col,
      };
    }

    $scope.removeWidget = function(index) {
      widgetsManager.remove($scope.widgets[index]);
      $scope.widgets.splice(index, 1);
    };

    $scope.addColumn = function() {
      $scope.gridsterOpts.columns += 1;
    };

    $scope.removeColumn = function() {
      if($scope.gridsterOpts.columns > 0) {
        $scope.gridsterOpts.columns -= 1;
      }
    };

    $scope.addWidgets = function() {
      widgetsManager.updateAvailableClients();
    };

    $scope.saveLayout = function() {
      var widgetsData = widgetsManager.serialize();
      var data = {
        columns: $scope.gridsterOpts.columns,
        widgets: widgetsData,

      };

      var pretty = angular.toJson(data, true);
      var urlEncoded = encodeURIComponent(angular.toJson(data));

      $log.debug('Layout data:');
      $log.debug(data);
      $log.debug(urlEncoded);

      $scope.saveLayout.raw = pretty;

      var url = '#?layout=' + urlEncoded;
      $scope.saveLayout.url = url;
    };

    $scope.goToLayoutURL = function() {
      $window.location.assign($scope.saveLayout.url);
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

    $scope.widgets = [];

    $rootScope.$on(widgetsManager.addWidgetRequestEvent, function(ev, client, type) {
      var w = widgetsManager.add(client, type);
      w.client = client;

      var position = findFreeTile(w.width);
      w.col = position.column;
      w.row = position.row;

      $scope.widgets.push(w);
      widgetsManager.update(true);
    });

    widgetsManager.start();

    var query = $location.search();
    if("layout" in query) {
      try {
        var json = angular.fromJson(query.layout);
        var columns = json.columns;
        var loaded = widgetsManager.deserialize(json.widgets);

        $scope.gridsterOpts.columns = columns;
        $scope.widgets = loaded;

      } catch(error) {
        $log.error('Could not load layout: invalid JSON:', error);
      }
    }

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
      if(!$scope.clients || !$scope.clients.length) {
        $scope.message = "No clients connected. Connect a client and refresh.";
      } else {
        $scope.message = null;
      }
    }

    update();
    $rootScope.$on(widgetsManager.availableClientsChangedEvent, function() {
      update();
    });

    $scope.addWidget = function(clientIdentifier, widgetType) {
      $rootScope.$emit(widgetsManager.addWidgetRequestEvent, clientIdentifier, widgetType);
    };

    $scope.refresh = function() {
      widgetsManager.updateAvailableClients();
    };
  }
]);
