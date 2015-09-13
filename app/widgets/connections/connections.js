'use strict';

angular.module('dashboard').controller('ConnectionsWidgetController', [
  function() {
  }
]);

angular.module('dashboard').factory('ConnectionsWidget', [
  'Widget',
  function(Widget) {
    return function() {
      var self = new Widget('connections-widget', 'Connections');

      self.data = {
        tcp4: 0,
        tcp6: 0,
      };

      return self;
    };
  }
]);

angular.module('dashboard').directive('connectionsWidget', [
  function() {
    return {
      replace: true,
      templateUrl: 'widgets/connections/connections.html',
    };
  }
]);
