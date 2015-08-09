'use strict';

angular.module('dashboard').controller('GridController', [
  '$scope',
  'widgetsManager',
  'Widget',
  function($scope, widgetsManager, Widget) {
    widgetsManager.register(new Widget('webserver', 'load', {
      sizeX: 1, sizeY: 1, row: 0, col: 0,
    }));
  }
]);
