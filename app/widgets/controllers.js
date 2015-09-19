'use strict';

angular.module('dashboard').controller('GridController', [
  '$scope',
  '$rootScope',
  '$window',
  '$location',
  'widgetsManager',
  'widgetsStore',
  'LoadWidget',
  'UptimeWidget',
  '$log',
  function($scope, $rootScope, $window, $location, widgetsManager, widgetsStore, LoadWidget, UptimeWidget, $log) {
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

    $scope.hoverMenu = false;

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

    $scope.clearWidgets = function() {
      widgetsStore.clearLayout();
      $window.location.href = '/';
    };

    $scope.saveLayout = function() {
      var data = widgetsStore.serialize();
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
    widgetsStore.columns = 4;
    $scope.$watch('gridsterOpts.columns', function() {
      widgetsStore.columns = $scope.gridsterOpts.columns;
      widgetsStore.saveLayout();
    });

    $scope.widgets = [];
    $scope.$watch(function() {
      var converted = [];
      $scope.widgets.forEach(function(widget) {
        converted.push(widget.watchValue());
      });
      return converted;
    }, function(widgets) {
      widgetsStore.widgets = widgets;
      widgetsStore.saveLayout();
    }, true);

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

    var data = widgetsStore.loadLayout();
    if(data) {
      try {
        $log.debug('Using', data);
        var json = angular.fromJson(data);
        $log.debug('Got:', json);
        var columns = json.columns;
        var loaded = widgetsManager.addFrom(json.widgets);

        $scope.gridsterOpts.columns = columns;
        $scope.widgets = loaded;

        widgetsStore.saveLayout();

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
      var clients = widgetsManager.availableClients;
      $log.debug('Available:', clients);

      var allClients = [
        {
          name: 'Built-In',
          identifier: '',
          availableWidgets : [
            'current_time',
            'current_date',
          ],
        },
      ];

      $scope.clients = allClients.concat(clients);
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
