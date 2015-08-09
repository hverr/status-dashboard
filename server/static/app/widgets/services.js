'use strict';

angular.module('dashboard').factory('widgetsManager', [
  '$log',
  function($log) {
    var widgets = [];

    return {
      register : function(widget) {
        $log.info('Registering widget', widget.type, 'for', widget.client);
        widgets.push(widget);
      },

      registeredWidgets : function() {
        return widgets;
      },
    };
  }
]);
