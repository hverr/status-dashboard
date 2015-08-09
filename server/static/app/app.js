'use strict';

// Declare app level module which depends on views, and components
angular.module('dashboard', [
  'ngRoute',
  'gridster',
]);

angular.module('dashboard').run([
  'widgetsManager',
  function(widgetsManager) {
    widgetsManager.start();
  }
]);
