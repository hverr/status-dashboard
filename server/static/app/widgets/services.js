'use strict';

angular.module('dashboard').factory('widgetsManager', [
  '$timeout',
  '$log',
  function($timeout, $log) {
    var widgets = [];

    return {
      register : function(widget) {
        $log.info('Registering widget', widget.type, 'for', widget.client);
        widgets.push(widget);
      },

      registeredWidgets : function() {
        return widgets;
      },

      start : function() {
        var f = function() {
          widgets.forEach(function(w) {
            w.update();
          });
        };
        $timeout(f, 3000);
      },
    };
  }
]);
