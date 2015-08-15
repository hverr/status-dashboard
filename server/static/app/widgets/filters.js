'use strict';

angular.module('dashboard').filter('widgetName', [
  'widgetFactory',
  function(widgetFactory) {
    return function(input) {
      var w = widgetFactory(input);
      if(w === null) {
        return 'Unknown';
      }

      return w.name;
    };
  }
]);
