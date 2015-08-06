'use strict';

// Declare app level module which depends on views, and components
angular.module('dashboard', [
  'ngRoute',
  'gridster',
])
.controller('gridController', [
  '$scope',
  function($scope) {
    $scope.standardItems = [
      { sizeX: 2, sizeY: 1, row: 0, col: 0 },
      { sizeX: 2, sizeY: 2, row: 0, col: 2 },
      { sizeX: 1, sizeY: 1, row: 0, col: 4 },
      { sizeX: 2, sizeY: 1, row: 1, col: 0 },
      { sizeX: 1, sizeY: 1, row: 1, col: 4 },
      { sizeX: 5, sizeY: 2, row: 2, col: 0 },
    ];
  }
]);
